package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kentik/common/cmetrics/httptsdb"
	"github.com/kentik/go-metrics"
	"github.com/kentik/libkflow/agg"
)

const (
	MaxHttpRequests    = 3
	MetricsSampleSize  = 1028
	MetricsSampleAlpha = 0.015
)

type Metrics struct {
	agg.Metrics
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
		agg.Metrics{
			TotalFlowsIn:   metrics.GetOrRegisterMeter(name("Total"), nil),
			TotalFlowsOut:  metrics.GetOrRegisterMeter(name("TotalFlowsOut"), nil),
			OrigSampleRate: metrics.GetOrRegisterHistogram(name("OrigSampleRate"), nil, sample()),
			NewSampleRate:  metrics.GetOrRegisterHistogram(name("NewSampleRate"), nil, sample()),
			RateLimitDrops: metrics.GetOrRegisterMeter(name("RateLimitDrops"), nil),
		},
	}
}

func (m *Metrics) Start(url, email, token string, interval time.Duration) {
	extra := map[string]string{
		"ver":   "libkflow-" + Version,
		"ft":    "nprobe",
		"dt":    "libkflow",
		"level": "primary",
	}

	go httptsdb.OpenTSDBWithConfig(httptsdb.OpenTSDBConfig{
		Addr:               url,
		Registry:           metrics.DefaultRegistry,
		FlushInterval:      interval,
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
