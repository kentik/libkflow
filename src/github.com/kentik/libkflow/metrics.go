package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kentik/common/cmetrics/httptsdb"
	"github.com/kentik/go-metrics"
)

const (
	MaxHttpRequests    = 3
	MetricsSampleSize  = 1028
	MetricsSampleAlpha = 0.015
)

type Metrics struct {
	Total            metrics.Meter
	Pkts             metrics.Meter
	TotalFlowsOut    metrics.Meter
	ReadBacklog      metrics.Meter
	WriteBacklog     metrics.Meter
	DownsampleFPS    metrics.Meter
	OrigSampleRate   metrics.Histogram
	NewSampleRate    metrics.Histogram
	DiscardedPackets metrics.Meter // FIXME: don't need?
	RateLimitDrops   metrics.Meter
}

func NewMetrics(clientid string) *Metrics {
	clientid = strings.Replace(clientid, ":", ".", -1)

	name := func(key string) string {
		return fmt.Sprintf("client_%s.%s", key, clientid)
	}

	sample := func() metrics.Sample {
		return metrics.NewExpDecaySample(MetricsSampleSize, MetricsSampleAlpha)
	}

	return &Metrics{
		Total:            metrics.GetOrRegisterMeter(name("Total"), nil),
		Pkts:             metrics.GetOrRegisterMeter(name("Pkts"), nil),
		ReadBacklog:      metrics.GetOrRegisterMeter(name("ReadBacklog"), nil),
		WriteBacklog:     metrics.GetOrRegisterMeter(name("WriteBacklog"), nil),
		DownsampleFPS:    metrics.GetOrRegisterMeter(name("DownsampleFPS"), nil),
		OrigSampleRate:   metrics.GetOrRegisterHistogram(name("OrigSampleRate"), nil, sample()),
		NewSampleRate:    metrics.GetOrRegisterHistogram(name("NewSampleRate"), nil, sample()),
		DiscardedPackets: metrics.GetOrRegisterMeter(name("DiscardedPackets"), nil),
		RateLimitDrops:   metrics.GetOrRegisterMeter(name("RateLimitDrops"), nil),
		TotalFlowsOut:    metrics.GetOrRegisterMeter(name("TotalFlowsOut"), nil),
	}
}

func (m *Metrics) Start(url, email, token string) {
	extra := map[string]string{
		"ver":   "libkflow-0.0.0", // FIXME: proper version
		"ft":    "nprobe",
		"dt":    "libkflow",
		"level": "primary",
	}

	go httptsdb.OpenTSDBWithConfig(httptsdb.OpenTSDBConfig{
		Addr:               url,
		Registry:           metrics.DefaultRegistry,
		FlushInterval:      1 * time.Second,
		DurationUnit:       time.Millisecond,
		Prefix:             "chf",
		Debug:              false,
		Send:               make(chan []byte, MaxHttpRequests),
		ProxyUrl:           os.Getenv("CH_HTTP_LOCAL_PROXY"),
		MaxHttpOutstanding: MaxHttpRequests,
		Extra:              extra,
		ApiEmail:           &email,
		ApiPassword:        &token,
	})
}

func (m *Metrics) Update(packets uint64, backlog int) {
	if m == nil {
		return
	}

}
