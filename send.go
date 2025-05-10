package libkflow

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/tinylib/msgp/msgp"
	capnp "zombiezen.com/go/capnproto2"

	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/log"
	"github.com/kentik/libkflow/metrics"
)

var messagePrefix = [80]byte{}

// A Sender aggregates and transmits flow information to Kentik.
type Sender struct {
	agg               *agg.Agg
	exit              chan struct{}
	url               *url.URL
	timeout           time.Duration
	client            *api.Client
	sample            int
	ticker            *time.Ticker
	tickerCtx         context.Context
	tickerCancelFunc  context.CancelFunc
	workers           sync.WaitGroup
	dns               chan []byte
	Device            *api.Device
	Errors            chan<- error
	useInternalErrors bool
	Metrics           *metrics.Metrics
}

func newSender(url *url.URL, timeout time.Duration) *Sender {
	tickerCtx, cancelFunc := context.WithCancel(context.Background())
	return &Sender{
		exit:             make(chan struct{}),
		url:              url,
		timeout:          timeout,
		ticker:           time.NewTicker(20 * time.Minute),
		tickerCtx:        tickerCtx,
		tickerCancelFunc: cancelFunc,
	}
}

// Send adds a flow record to the outgoing queue.
func (s *Sender) Send(flow *flow.Flow) {
	log.Debugf("sending flow to aggregator")
	flow.DeviceId = uint32(s.Device.ID)
	s.agg.Add(flow)
}

// SendFlows sends the flows to the Kentik API, returning the number of bytes sent as the payload.
func (s *Sender) SendFlows(flows []flow.Flow) (int64, error) {
	if s.Device == nil {
		return 0, fmt.Errorf("device not initialized")
	}
	if len(flows) == 0 {
		return 0, nil
	}

	if s.Metrics != nil {
		s.Metrics.TotalFlowsIn.Mark(int64(len(flows)))
	}

	// ensure all flows have the device ID set; otherwise it may not be properly queried
	for i := range flows {
		flows[i].DeviceId = uint32(s.Device.ID)
	}

	// serialize the data
	_, segment, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return 0, fmt.Errorf("failed to create capn proto segment: %w", err)
	}
	message, err := flow.ToCapnProtoMessage(flows, segment)
	if err != nil {
		return 0, fmt.Errorf("failed to convert flows to capn proto: %w", err)
	}

	// write the data with additional gzip compression
	buf := &bytes.Buffer{}
	z := gzip.NewWriter(buf)
	_, err = z.Write(messagePrefix[:])
	if err != nil {
		return 0, fmt.Errorf("failed to write empty message header: %w", err)
	}
	err = capnp.NewPackedEncoder(z).Encode(message)
	if err != nil {
		return 0, fmt.Errorf("failed to encode packed capn proto message: %w", err)
	}
	err = z.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// send the compressed and packed message to the Kentik API
	payloadLength := int64(len(buf.Bytes()))
	err = s.client.SendFlow(s.url.String(), buf)
	if err != nil {
		return 0, err
	}

	if s.Metrics != nil {
		s.Metrics.TotalFlowsOut.Mark(int64(len(flows)))
		s.Metrics.BytesSent.Mark(payloadLength)
	}

	return payloadLength, nil
}

// Stop requests a graceful shutdown of the Sender.
func (s *Sender) Stop(wait time.Duration) bool {
	s.agg.Stop()
	select {
	case <-s.exit:
		return true
	case <-time.After(wait):
		return false
	}
}

func (s *Sender) GetClient() *api.Client {
	if s != nil {
		return s.client
	} else {
		return nil
	}
}

func (s *Sender) GetDevice() *api.Device {
	if s != nil {
		return s.Device
	}
	return nil
}

func (s *Sender) StartDNS(url *url.URL, interval time.Duration) {
	s.dns = make(chan []byte, 1e5)
	go s.dispatchDNS(url.String(), interval)
}

func (s *Sender) SendDNS(res *api.DNSResponse) error {
	buf := bytes.Buffer{}
	enc := msgp.NewWriter(&buf)

	err := res.EncodeMsg(enc)
	if err != nil {
		return err
	}

	enc.Flush()
	s.dns <- buf.Bytes()

	return nil
}

func (s *Sender) SendEncodedDNS(data []byte) {
	s.dns <- data
}

func (s *Sender) start(agg *agg.Agg, client *api.Client, device *api.Device, n int) error {
	q := s.url.Query()
	q.Set("sid", "0")
	q.Set("sender_id", device.ClientID())

	s.agg = agg
	s.url.RawQuery = q.Encode()
	s.Device = device
	s.client = client
	s.workers.Add(n)

	for i := 0; i < n; i++ {
		go s.dispatch()
	}
	go s.monitor()
	go s.update()

	log.Debugf("sender started with %d workers", n)

	return nil
}

func (s *Sender) dispatch() {
	buf := &bytes.Buffer{}
	cid := [80]byte{}
	url := s.url.String()
	z := gzip.NewWriter(buf)

	for msg := range s.agg.Output() {
		log.Debugf("dispatching aggregated flow")
		z.Reset(buf)
		z.Write(cid[:])

		err := capnp.NewPackedEncoder(z).Encode(msg)
		if err != nil {
			s.error(err)
			continue
		}

		z.Close()
		l := buf.Len()
		err = s.client.SendFlow(url, buf)
		if err != nil {
			s.error(err)
			continue
		}

		if s.Metrics != nil {
			s.Metrics.BytesSent.Mark(int64(l))
		}
	}
	s.workers.Done()
}

func (s *Sender) dispatchDNS(url string, interval time.Duration) {
	ticker := time.NewTicker(interval)

	buf := bytes.Buffer{}

	for {
		flush := false

		select {
		case data := <-s.dns:
			buf.Write(data)
		case <-ticker.C:
			flush = true
		}

		if buf.Len() > 1e6 || flush && buf.Len() > 0 {
			err := s.client.SendDNS(url, &buf)
			if err != nil {
				s.error(err)
				continue
			}

			buf.Reset()
		}
	}
}

func (s *Sender) monitor() {
	for {
		select {
		case err := <-s.agg.Errors():
			s.error(err)
		case <-s.agg.Done():
			s.workers.Wait()
			s.ticker.Stop()
			s.tickerCancelFunc()
			s.Metrics.Unregister()
			if s.useInternalErrors {
				close(s.Errors)
			}
			s.exit <- struct{}{}
			log.Debugf("sender stopped")
			return
		}
	}
}

func (s *Sender) update() {
	for {
		select {
		case <-s.tickerCtx.Done():
			return

		case <-s.ticker.C:
			updated, err := s.client.GetDeviceByID(s.Device.ID)
			if err != nil {
				if api.IsErrorWithStatusCode(err, 404) {
					updated = &api.Device{}
				} else {
					log.Debugf("device API request failed: %s", err)
					continue
				}
			}

			if s.Device.MaxFlowRate != updated.MaxFlowRate {
				log.Debugf("updating max FPS to %d", updated.MaxFlowRate)
				s.Device.MaxFlowRate = updated.MaxFlowRate
				s.agg.Configure(updated.MaxFlowRate)
			}

			// if the configured sample rate is 0 then the sender
			// may be using the device sample rate which has just
			// changed, so abort the program
			if s.sample == 0 && s.Device.SampleRate != updated.SampleRate {
				log.Debugf("device sample rate changed, aborting")
				os.Exit(1)
			}
		}
	}
}

func (s *Sender) error(err error) {
	select {
	case s.Errors <- err:
	default:
	}
}
