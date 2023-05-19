package metrics

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/kentik/kit/go/legacy/common/cmetrics"
	"github.com/kentik/kit/go/legacy/go-metrics"
)

const (
	MaxHttpRequests    = 3
	MetricsSampleSize  = 1028
	MetricsSampleAlpha = 0.015
)

type Metrics struct {
	reg            metrics.Registry
	TotalFlowsIn   metrics.Meter
	TotalFlowsOut  metrics.Meter
	OrigSampleRate metrics.Histogram
	NewSampleRate  metrics.Histogram
	RateLimitDrops metrics.Meter
	BytesSent      metrics.Meter
	Extra          map[string]string
}

func New(companyID int, deviceID int, program, version string) *Metrics {
	name := func(key string) string {
		return fmt.Sprintf("client_%s", key)
	}

	sample := func() metrics.Sample {
		return metrics.NewExpDecaySample(MetricsSampleSize, MetricsSampleAlpha)
	}

	extra := map[string]string{
		"ver":   program + "-" + version,
		"ft":    program,
		"dt":    "libkflow",
		"level": "primary",
		"cid":   strconv.Itoa(companyID),
		"did":   strconv.Itoa(deviceID),
	}

	// libkflow creates its own go-metrics Registry, which hold only its
	// own metrics (or ones that its clients create with
	reg := metrics.NewRegistry()

	return &Metrics{
		reg:            reg,
		TotalFlowsIn:   metrics.GetOrRegisterMeter(name("Total"), reg),
		TotalFlowsOut:  metrics.GetOrRegisterMeter(name("DownsampleFPS"), reg),
		OrigSampleRate: metrics.GetOrRegisterHistogram(name("OrigSampleRate"), reg, sample()),
		NewSampleRate:  metrics.GetOrRegisterHistogram(name("NewSampleRate"), reg, sample()),
		RateLimitDrops: metrics.GetOrRegisterMeter(name("RateLimitDrops"), reg),
		BytesSent:      metrics.GetOrRegisterMeter(name("BytesSent"), reg),
		Extra:          extra,
	}
}

func (m *Metrics) Start(url, email, token string, interval time.Duration, proxy *url.URL) {
	proxyURL := ""
	if proxy != nil {
		proxyURL = proxy.String()
	}

	go cmetrics.OpenHTTPTSDBWithConfig(cmetrics.OpenHTTPTSDBConfig{
		Addr:               url,
		Registry:           m.reg,
		FlushInterval:      interval,
		DurationUnit:       time.Millisecond,
		Prefix:             "chf",
		Debug:              false,
		Send:               make(chan []byte, MaxHttpRequests),
		ProxyUrl:           proxyURL,
		MaxHttpOutstanding: MaxHttpRequests,
		Extra:              m.Extra,
		ApiEmail:           &email,
		ApiPassword:        &token,
	})
}

func NewMeter() metrics.Meter {
	return metrics.NewMeter()
}

func NewHistogram(s metrics.Sample) metrics.Histogram {
	return metrics.NewHistogram(s)
}

func NewUniformSample(n int) metrics.Sample {
	return metrics.NewUniformSample(n)
}
