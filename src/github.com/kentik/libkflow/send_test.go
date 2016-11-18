package main

import (
	"io/ioutil"
	"math/rand"
	"net/url"
	"testing"
	"time"

	"github.com/kentik/go-metrics"
	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api/test"
	"github.com/kentik/libkflow/chf"
	"github.com/stretchr/testify/assert"
)

func TestSender(t *testing.T) {
	sender, server, assert := setup(t)

	expected, err := chf.NewCHF(sender.Segment())
	if err != nil {
		t.Fatal(err)
	}
	expected.SetSrcAs(rand.Uint32())
	expected.SetDstAs(rand.Uint32())

	sender.Send(&expected)
	msgs, err := (<-server.Flows()).Msgs()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expected.String(), msgs.At(0).String())
}

func TestSenderStop(t *testing.T) {
	sender, _, assert := setup(t)
	stopped := sender.Stop(100 * time.Millisecond)
	assert.True(stopped)
}

func setup(t *testing.T) (*Sender, *test.Server, *assert.Assertions) {
	metrics := &agg.Metrics{
		TotalFlowsIn:   metrics.NewMeter(),
		TotalFlowsOut:  metrics.NewMeter(),
		OrigSampleRate: metrics.NewHistogram(metrics.NewUniformSample(100)),
		NewSampleRate:  metrics.NewHistogram(metrics.NewUniformSample(100)),
		RateLimitDrops: metrics.NewMeter(),
	}

	agg, err := agg.NewAgg(10*time.Millisecond, 10, metrics)
	if err != nil {
		t.Fatal(err)
	}

	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	server.Log.SetOutput(ioutil.Discard)

	url, _ := url.Parse(server.URL() + "/chf")
	sender := NewSender(url, 1*time.Second, 1)
	sender.Start(agg, client, device, 1)

	return sender, server, assert.New(t)
}
