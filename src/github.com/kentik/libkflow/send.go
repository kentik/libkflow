package main

import (
	"bytes"
	"log"
	"net/url"
	"time"

	"github.com/kentik/libkflow/api"
	"zombiezen.com/go/capnproto2"
)

type Sender struct {
	Flow    chan *capnp.Message
	Exit    chan struct{}
	URL     *url.URL
	Timeout time.Duration
	Client  *api.Client
	Verbose int
}

func NewSender(url *url.URL, timeout time.Duration, verbose int) *Sender {
	return &Sender{
		Flow:    make(chan *capnp.Message, 100),
		Exit:    make(chan struct{}),
		URL:     url,
		Timeout: timeout,
		Verbose: verbose,
	}
}

func (s *Sender) Validate(url, email, token string, did int) error {
	s.Client = api.NewClient(email, token, s.Timeout)

	device, err := s.Client.GetDevice(url, did)
	if err != nil {
		return err
	}

	// FIXME: use custom columns

	q := s.URL.Query()
	q.Set("sid", "0")
	q.Set("sender_id", device.ClientID())
	s.URL.RawQuery = q.Encode()

	s.Client.Header.Set("Content-Type", "application/binary")

	go s.dispatch()

	return nil
}

func (s *Sender) Send(msg *capnp.Message) bool {
	select {
	case s.Flow <- msg:
		return true
	default:
		return false
	}
}

func (s *Sender) Stop(wait time.Duration) bool {
	s.Exit <- struct{}{}
	select {
	case <-s.Exit:
		return true
	case <-time.After(wait):
		return false
	}
}

func (s *Sender) dispatch() {
	buf := &bytes.Buffer{}
	cid := [80]byte{}
	url := s.URL.String()

	for {
		select {
		case msg := <-s.Flow:
			buf.Reset()
			buf.Write(cid[:])

			err := capnp.NewPackedEncoder(buf).Encode(msg)
			if err != nil {
				// FIXME: check verbosity
				log.Print("NewPackedEncoder", err)
				continue
			}

			err = s.Client.SendFlow(url, buf)
			if err != nil {
				// FIXME: check verbosity
				log.Print("HTTP", err)
				continue
			}
		case <-s.Exit:
			s.Exit <- struct{}{}
			return
		}
	}
}
