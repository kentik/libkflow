package libkflow_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kentik/kit/go/legacy/go-metrics"
	"github.com/kentik/libkflow"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
	"github.com/kentik/libkflow/flow"
	metrics2 "github.com/kentik/libkflow/metrics"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestNewSenderWithDeviceID(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderWithDeviceID(dev.ID, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func TestNewSenderWithDeviceIP(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderWithDeviceIP(dev.IP, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func TestNewSenderWithDeviceName(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderWithDeviceName(dev.Name, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func TestNewSenderWithDeviceNameWithErrors_NoErrs(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	apiurl = server.URL(test.API)
	flowurl = server.URL(test.FLOW)
	metricsurl = server.URL(test.TSDB)

	email = client.Email
	token = client.Token

	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	l := stubLeveledLogger{}

	registry := metrics.NewRegistry()
	metrics2.StartWithSetConf(registry, &l, metricsurl.String(), email, token, "chf")
	config.WithRegistry(registry)

	s, errors, err := libkflow.NewSenderWithDeviceNameWithErrors(device.Name, config)
	assert.NoError(t, err)

	errorsFromChan := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for err := range errors {
			errorsFromChan = append(errorsFromChan, err)
		}
		wg.Done()
	}()

	for i := 0; i < 5; i++ {
		s.Send(&flow.Flow{
			TimestampNano: time.Now().UnixNano(),
		})
	}

	s.Stop(time.Second)

	wg.Wait()

	assert.Len(t, errorsFromChan, 0)
}

func TestNewSenderWithDeviceNameWithErrors_WithErrs(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	flowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))

	apiurl = server.URL(test.API)
	flowurl = server.URL(flowServer.URL)
	metricsurl = server.URL(test.TSDB)

	email = client.Email
	token = client.Token

	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	l := stubLeveledLogger{}

	registry := metrics.NewRegistry()
	metrics2.StartWithSetConf(registry, &l, metricsurl.String(), email, token, "chf")
	config.WithRegistry(registry)

	s, errors, err := libkflow.NewSenderWithDeviceNameWithErrors(device.Name, config)
	assert.NoError(t, err)

	errorsFromChan := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for err := range errors {
			errorsFromChan = append(errorsFromChan, err)
		}
		wg.Done()
	}()

	s.Send(&flow.Flow{
		TimestampNano: time.Now().UnixNano(),
	})

	s.Stop(time.Second)

	wg.Wait()

	assert.Len(t, errorsFromChan, 1)
}

func TestNewSenderWithDeviceName_WithErrs_NoPanic(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	flowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		time.Sleep(time.Second)
	}))

	apiurl = server.URL(test.API)
	flowurl = server.URL(flowServer.URL)
	metricsurl = server.URL(test.TSDB)

	email = client.Email
	token = client.Token

	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	l := stubLeveledLogger{}

	registry := metrics.NewRegistry()
	metrics2.StartWithSetConf(registry, &l, metricsurl.String(), email, token, "chf")
	config.WithRegistry(registry)

	errors := make(chan error)

	s, err := libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(t, err)

	errorsFromChan := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Second):
				return
			case err := <-errors:
				errorsFromChan = append(errorsFromChan, err)
			}
		}
	}()

	for i := 0; i < 100000; i++ {
		s.Send(&flow.Flow{
			TimestampNano: time.Now().UnixNano(),
		})
	}

	s.Stop(time.Second * 0)

	wg.Wait()

	assert.Len(t, errorsFromChan, 1)
}

func TestNewSenderWithDeviceNameLeaks(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	apiurl = server.URL(test.API)
	flowurl = server.URL(test.FLOW)
	metricsurl = server.URL(test.TSDB)

	email = client.Email
	token = client.Token

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	l := stubLeveledLogger{}

	registry := metrics.NewRegistry()
	metrics2.StartWithSetConf(registry, &l, metricsurl.String(), email, token, "chf")
	config.WithRegistry(registry)

	// kick off the tick go routines from the go metrics library
	// these are only started once per process
	_ = metrics.NewMeter()

	ignore := goleak.IgnoreCurrent()

	s, err := libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(err)
	assert.NotNil(s)
	s.Stop(time.Second)

	s, err = libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(err)
	assert.NotNil(s)
	s.Stop(time.Second)

	s, err = libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(err)
	assert.NotNil(s)
	s.Stop(time.Second)

	s, err = libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(err)
	assert.NotNil(s)
	s.Stop(time.Second)

	s, err = libkflow.NewSenderWithDeviceName(device.Name, errors, config)
	assert.NoError(err)
	assert.NotNil(s)
	s.Stop(time.Second)

	server.Close()

	time.Sleep(time.Second)

	goleak.VerifyNone(t, ignore)
}

func TestNewSenderFromDevice(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token, "test", "0.0.1")
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderFromDevice(dev, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func setupLibTest(t *testing.T) (*api.Device, *assert.Assertions) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	apiurl = server.URL(test.API)
	flowurl = server.URL(test.FLOW)
	metricsurl = server.URL(test.TSDB)

	email = client.Email
	token = client.Token

	return device, assert
}

var (
	apiurl     *url.URL
	flowurl    *url.URL
	metricsurl *url.URL
	email      string
	token      string
)

type stubLogger struct {
	count uint32
}

func (s *stubLogger) Printf(string, ...interface{}) { atomic.AddUint32(&s.count, 1) }

type stubLeveledLogger struct {
	count uint32
}

func (s *stubLeveledLogger) Errorf(string, string, ...interface{}) { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Infof(string, string, ...interface{})  { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Debugf(string, string, ...interface{}) { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Warnf(string, string, ...interface{})  { atomic.AddUint32(&s.count, 1) }
