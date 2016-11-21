package main

import (
	"bytes"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/chf"
	"zombiezen.com/go/capnproto2"
)

type Sender struct {
	Agg     *agg.Agg
	Exit    chan struct{}
	URL     *url.URL
	Timeout time.Duration
	Client  *api.Client
	Verbose int
	Customs api.CustomColumns
	workers sync.WaitGroup
}

func NewSender(url *url.URL, timeout time.Duration, verbose int) *Sender {
	return &Sender{
		Exit:    make(chan struct{}),
		URL:     url,
		Timeout: timeout,
		Verbose: verbose,
	}
}

func (s *Sender) Start(agg *agg.Agg, client *api.Client, device *api.Device, n int) error {
	client.Header.Set("Content-Type", "application/binary")

	q := s.URL.Query()
	q.Set("sid", "0")
	q.Set("sender_id", device.ClientID())

	s.Agg = agg
	s.URL.RawQuery = q.Encode()
	s.Customs = device.Customs
	s.Client = client
	s.workers.Add(n)

	for i := 0; i < n; i++ {
		go s.dispatch()
	}
	go s.monitor()

	s.debug("sender started with %d workers", n)

	return nil
}

func (s *Sender) Segment() *capnp.Segment {
	return s.Agg.Segment()
}

func (s *Sender) Send(flow *chf.CHF) {
	s.debug("sending flow to aggregator")
	s.Agg.Add(flow)
}

func (s *Sender) Stop(wait time.Duration) bool {
	s.Agg.Stop()
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

	for msg := range s.Agg.Output() {
		s.debug("dispatching aggregated flow")

		buf.Reset()
		buf.Write(cid[:])

		err := capnp.NewPackedEncoder(buf).Encode(msg)
		if err != nil {
			log.Print(err)
			continue
		}

		err = s.Client.SendFlow(url, buf)
		if err != nil {
			log.Print(err)
			continue
		}
	}
	s.workers.Done()
}

func (s *Sender) monitor() {
	for {
		select {
		case err := <-s.Agg.Errors():
			log.Print(err)
		case <-s.Agg.Done():
			s.workers.Wait()
			s.Exit <- struct{}{}
			s.debug("sender stopped")
			return
		}
	}
}

func (s *Sender) debug(fmt string, v ...interface{}) {
	if s.Verbose > 0 {
		log.Printf(fmt, v...)
	}
}
