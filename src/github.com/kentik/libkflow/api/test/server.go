package test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/chf"
	"zombiezen.com/go/capnproto2"
)

type Server struct {
	Host     net.IP
	Port     int
	Email    string
	Token    string
	Device   api.Device
	Log      *log.Logger
	flows    chan chf.PackedCHF
	mux      *http.ServeMux
	listener net.Listener
}

func NewServer(host string, port int, tls bool) (*Server, error) {
	var listener net.Listener

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	if tls {
		listener, err = tlslistener(listener, host, addr)
		if err != nil {
			return nil, err
		}
	}

	addr = listener.Addr().(*net.TCPAddr)

	return &Server{
		Host:     addr.IP,
		Port:     addr.Port,
		Log:      log.New(os.Stderr, "", log.LstdFlags),
		flows:    make(chan chf.PackedCHF, 10),
		mux:      http.NewServeMux(),
		listener: listener,
	}, nil
}

func (s *Server) Serve(email, token string, dev api.Device) error {
	s.Email = email
	s.Token = token
	s.Device = dev
	s.mux.HandleFunc("/api/v5/device/", s.wrap(s.device))
	s.mux.HandleFunc("/chf", s.wrap(s.flow))
	s.mux.HandleFunc("/tsdb", s.wrap(s.tsdb))
	return http.Serve(s.listener, s.mux)
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s *Server) Flows() <-chan chf.PackedCHF {
	return s.flows
}

func (s *Server) device(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[4]

	if id != strconv.Itoa(s.Device.ID) && id != s.Device.Name {
		panic(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&api.DeviceResponse{
		Device: s.Device,
	})

	if err != nil {
		panic(http.StatusInternalServerError)
	}
}

func (s *Server) flow(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("sid") != "0" {
		panic(http.StatusBadRequest)
	}

	if r.FormValue("sender_id") != s.Device.ClientID() {
		panic(http.StatusBadRequest)
	}

	cid := [80]byte{}
	n, err := r.Body.Read(cid[:])
	if err != nil || n != len(cid) {
		panic(http.StatusBadRequest)
	}

	msg, err := capnp.NewPackedDecoder(r.Body).Decode()
	defer r.Body.Close()
	if err != nil {
		panic(http.StatusBadRequest)
	}

	root, err := chf.ReadRootPackedCHF(msg)
	if err != nil {
		panic(http.StatusBadRequest)
	}

	select {
	case s.flows <- root:
	default:
	}

	msgs, err := root.Msgs()
	if err != nil {
		panic(http.StatusBadRequest)
	}

	for i := 0; i < msgs.Len(); i++ {
		buf := bytes.Buffer{}
		Print(&buf, i, msgs.At(i))
		s.Log.Output(0, buf.String())
	}
}

func (s *Server) tsdb(w http.ResponseWriter, r *http.Request) {
	// just ignore it
}

func (s *Server) wrap(f handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if code, ok := r.(int); ok {
					http.Error(w, http.StatusText(code), code)
					return
				}
				panic(r)
			}
		}()

		email := r.Header.Get("X-CH-Auth-Email")
		token := r.Header.Get("X-CH-Auth-API-Token")

		if email != s.Email || token != s.Token {
			panic(http.StatusUnauthorized)
		}

		if err := r.ParseForm(); err != nil {
			panic(http.StatusBadRequest)
		}

		f(w, r)
	}
}

func tlslistener(tcp net.Listener, host string, addr *net.TCPAddr) (net.Listener, error) {
	pri, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := &pri.PublicKey

	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber:          sn,
		Subject:               pkix.Name{Organization: []string{"Kentik"}},
		IPAddresses:           []net.IP{addr.IP},
		DNSNames:              []string{host},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, pri)
	if err != nil {
		return nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  pri,
	}

	cfg := tls.Config{Certificates: []tls.Certificate{cert}}
	return tls.NewListener(tcp, &cfg), nil
}

type handler func(http.ResponseWriter, *http.Request)
