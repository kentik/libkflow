package main

import (
	"net/url"
	"testing"
	"time"

	"github.com/kentik/go-metrics"
	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

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

	agg, err := agg.NewAgg(100*time.Millisecond, 10, metrics)
	if err != nil {
		t.Fatal(err)
	}

	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	url, _ := url.Parse(server.URL() + "/chf")
	sender := NewSender(url, 1*time.Second, 1)
	sender.Start(agg, client, device)

	return sender, server, assert.New(t)
}
