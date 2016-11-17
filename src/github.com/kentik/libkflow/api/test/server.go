package test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

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
	mux      *http.ServeMux
	listener net.Listener
}

func NewServer(host string, port int) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	addr = listener.Addr().(*net.TCPAddr)

	return &Server{
		Host:     addr.IP,
		Port:     addr.Port,
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
	return http.Serve(s.listener, s.mux)
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s *Server) device(w http.ResponseWriter, r *http.Request) {
	var did int

	n, err := fmt.Sscanf(r.URL.Path, "/api/v5/device/%d", &did)
	if n != 1 || err != nil {
		panic(http.StatusBadRequest)
	}

	if did != s.Device.ID {
		panic(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&api.DeviceResponse{
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

	msgs, err := root.Msgs()
	if err != nil {
		panic(http.StatusBadRequest)
	}

	for i := 0; i < msgs.Len(); i++ {
		Print(i, msgs.At(i))
	}
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

type handler func(http.ResponseWriter, *http.Request)
