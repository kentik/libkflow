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
	RateLimitDrops metrics.Meter
	BytesSent      metrics.Meter
}

func New(companyID int, deviceID int, program, version string) *Metrics {
	// libkflow creates its own go-metrics Registry, which hold only its
	// own metrics (or ones that its clients create with
	reg := metrics.NewRegistry()
	return NewWithRegistry(reg, companyID, deviceID, program, version)
}

// NewWithRegistry returns a new Metrics but allows a specific registry to be used rather than creating a new one
func NewWithRegistry(reg metrics.Registry, companyID int, deviceID int, program, version string) *Metrics {
	suffix := fmt.Sprintf("^ver=%s^ft=%s^dt=%s^level=%s^cid=%s^did=%s", program+"-"+version, program, "libkflow", "primary", strconv.Itoa(companyID), strconv.Itoa(deviceID))

	return &Metrics{
		reg:            reg,
		TotalFlowsIn:   metrics.GetOrRegisterMeter("client_Total"+suffix, reg),
		TotalFlowsOut:  metrics.GetOrRegisterMeter("client_DownsampleFPS"+suffix, reg),
		RateLimitDrops: metrics.GetOrRegisterMeter("client_RateLimitDrops"+suffix, reg),
		BytesSent:      metrics.GetOrRegisterMeter("client_BytesSent"+suffix, reg),
	}
}

func StartWithSetConf(registry metrics.Registry, logger cmetrics.Logger, url, email, token string, prefix string) {
	cmetrics.SetConfWithRegistry(url, logger, prefix, "chf", nil, nil, &email, &token, registry)
}

func (m *Metrics) Start(url, email, token string, prefix string, interval time.Duration, proxy *url.URL) {
	proxyURL := ""
	if proxy != nil {
		proxyURL = proxy.String()
	}

	go cmetrics.OpenHTTPTSDBWithConfig(cmetrics.OpenHTTPTSDBConfig{
		Addr:               url,
		Registry:           m.reg,
		FlushInterval:      interval,
		DurationUnit:       time.Millisecond,
		Prefix:             prefix,
		Debug:              false,
		Send:               make(chan []byte, MaxHttpRequests),
		ProxyUrl:           proxyURL,
		MaxHttpOutstanding: MaxHttpRequests,
		ApiEmail:           &email,
		ApiPassword:        &token,
	})
}

func (m *Metrics) Unregister() {
	m.reg.Unregister("client_Total")
	m.reg.Unregister("client_DownsampleFPS")
	m.reg.Unregister("client_OrigSampleRate")
	m.reg.Unregister("client_NewSampleRate")
	m.reg.Unregister("client_RateLimitDrops")
	m.reg.Unregister("client_BytesSent")
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
