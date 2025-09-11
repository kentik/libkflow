package cmetrics

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kentik/kit/go/legacy/go-metrics"
	"github.com/kentik/kit/go/legacy/golog/logger"
)

const (
	influxMaxSendTries          = 2
	influxClientResponseTimeout = 5 * time.Second
	influxClientKeepAlive       = 60 * time.Second
	influxClientTLSTimeout      = 5 * time.Second
	influxContentType           = "application/x-www-form-urlencoded"
	influxPrefix                = "[influx metrics]"
	maxHttpBodyLen              = 1000 * 1000 // 1MB
)

var myLogger Logger = logger.New(logger.Levels.Info)

type INFLUXMetricSet struct {
	keepZero bool // should metrics with zero values be kept or dropped
	metrics  []*INFLUXMetric
}

type INFLUXMetric struct {
	Metric    string
	Timestamp int64
	Value     int64
	Type      string
	Tags      map[string]string
}

// OpenINFLUXConfig provides a container with configuration parameters for
// the OpenINFLUX exporter
type OpenINFLUXConfig struct {
	Addr               string                // Network address to connect to
	Registry           metrics.Registry      // Registry to be exported
	FlushInterval      time.Duration         // Flush interval
	DurationUnit       time.Duration         // Time conversion unit for durations
	Prefix             string                // Prefix to be prepended to metric names
	Debug              bool                  // write to stdout for debug
	Quiet              bool                  // silence all comms
	Tags               map[string]string     // add these tags to each metric writen
	Send               chan *INFLUXMetricSet // manage # of outstanding http requests here.
	MaxHttpOutstanding int
	ProxyUrl           string
	Extra              map[string]string
	Args               map[string]string // extra config
	ApiEmail           *string
	ApiPassword        *string
}

// OpenINFLUX is a blocking exporter function which reports metrics in r
// to a INFLUX server located at addr, flushing them every d duration
// and prepending metric names with prefix.
func OpenINFLUX(r metrics.Registry, d time.Duration, prefix string, addr string, maxOutstanding int) {
	OpenINFLUXWithConfig(OpenINFLUXConfig{
		Addr:               addr,
		Registry:           r,
		FlushInterval:      d,
		DurationUnit:       time.Nanosecond,
		Prefix:             prefix,
		Debug:              false,
		MaxHttpOutstanding: maxOutstanding,
		Send:               make(chan *INFLUXMetricSet, maxOutstanding),
		Tags:               map[string]string{},
		Extra:              nil,
	})
}

// OpenINFLUXWithConfig is a blocking exporter function just like OpenINFLUX,
// but it takes a OpenINFLUXConfig instead.
func OpenINFLUXWithConfig(c OpenINFLUXConfig) {
	go c.runSend()

	// wait a bit for metrics to be registered
	time.Sleep(time.Second * 5)
	// and now do an initial export
	if err := openINFLUX(&c); nil != err {
		log.Println(err)
	}
	// then once every FlushInterval
	for range time.Tick(c.FlushInterval) {
		if err := openINFLUX(&c); nil != err {
			log.Println(err)
		}
	}
}

func (c *OpenINFLUXConfig) runSend() {
	if strings.HasPrefix(c.Addr, "http") {
		c.runSendViaHTTP()
	} else if strings.HasPrefix(c.Addr, "tcp") || strings.HasPrefix(c.Addr, "udp") {
		c.runSendViaSocket()
	}
}

func (c *OpenINFLUXConfig) runSendViaSocket() {
	for r := range c.Send {
		var w *bufio.Writer
		var conn net.Conn = nil

		if c.Debug {
			w = bufio.NewWriter(os.Stdout)
		} else {
			var err error
			pts := strings.Split(c.Addr, "://")
			if len(pts) == 2 {
				conn, err = net.Dial(pts[0], pts[1])
				if nil != err {
					if !c.Quiet {
						fmt.Printf("Invalid metrics address: %s, %v\n", c.Addr, err)
					}
					continue
				}
				w = bufio.NewWriter(conn)
			} else {
				if !c.Quiet {
					fmt.Printf("Invalid metrics address: %s\n", c.Addr)
				}
				continue
			}
		}

		if ebytes, err := r.ToWire(); err != nil {
			if !c.Quiet {
				fmt.Printf("Error encoding to wire format: %v\n", err)
			}
			continue
		} else {
			_, _ = w.Write(ebytes)
			_ = w.Flush()
		}

		if !c.Debug {
			_ = conn.Close()
		}
	}
}

func (c *OpenINFLUXConfig) runSendViaHTTP() {
	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		Dial: (&net.Dialer{
			Timeout:   influxClientResponseTimeout,
			KeepAlive: influxClientKeepAlive,
		}).Dial,
		TLSHandshakeTimeout: influxClientTLSTimeout,
	}

	// Add a proxy if needed.
	if c.ProxyUrl != "" {
		proxyUrl, err := url.Parse(c.ProxyUrl)
		if err != nil {
			myLogger.Errorf(influxPrefix, "ignored: error setting proxy: %s", err)
		} else {
			tr.Proxy = http.ProxyURL(proxyUrl)
			myLogger.Debugf(influxPrefix, "set outbound proxy: %s", c.ProxyUrl)
		}
	}

	for r := range c.Send {
		chunks, err := r.toWireChunks(maxHttpBodyLen)
		if err != nil {
			myLogger.Errorf(influxPrefix, "error encoding to wire format: %s", err)
			continue
		}
		for _, chunk := range chunks {
			if c.Debug {
				myLogger.Infof(influxPrefix, "metrics: %s", string(chunk))
				continue
			}

			for i := 0; i < influxMaxSendTries; i++ {
				req, err := http.NewRequest("POST", c.Addr, bytes.NewBuffer(chunk))
				if err != nil {
					myLogger.Errorf(influxPrefix, "error creating request: %s", err)
					continue
				}
				req.Header.Add("Content-Type", influxContentType)

				if c.ApiEmail != nil && c.ApiPassword != nil {
					req.Header.Add(API_EMAIL_HEADER, *c.ApiEmail)
					req.Header.Add(API_PASSWORD_HEADER, *c.ApiPassword)
				}

				client := &http.Client{Transport: tr, Timeout: influxClientResponseTimeout}
				resp, err := client.Do(req)
				if err != nil {
					if i == influxMaxSendTries-1 {
						myLogger.Errorf(influxPrefix, "posting to %s failed, giving up: %s", c.Addr, err)
					} else {
						myLogger.Warnf(influxPrefix, "retrying post to %s: %s", c.Addr, err)
					}
				} else {
					// Fire and forget
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()
					break
				}
			}

		}
	}
}

/*
*
Write out additional tags
*/
func openINFLUX(c *OpenINFLUXConfig) error {

	shortHostnameBase := GetShortHostname()
	now := time.Now().UnixNano()
	keepZero := c.Args["keep_zero"] == "true"
	sendBody := NewINFLUXMetricSet(keepZero)
	du := float64(c.DurationUnit)

	c.Registry.Each(func(baseName string, i interface{}) {
		name, tags := ExpandTags(baseName, c.Prefix, shortHostnameBase, c.Tags, c.Extra)

		switch metric := i.(type) {
		case metrics.Counter:
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: metric.Count(), Tags: tags, Type: "count"})
		case metrics.Gauge:
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: metric.Value(), Tags: tags, Type: "value"})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Min(), Tags: tags, Type: "min"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: h.Max(), Tags: tags, Type: "max"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(h.Mean()), Tags: tags, Type: "mean"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[2]), Tags: tags, Type: "95-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[3]), Tags: tags, Type: "99-percentile"})
			metric.Clear()
		case metrics.Meter:
			m := metric.Snapshot()
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: m.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(m.Rate1()), Tags: tags, Type: "one-minute"})
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Count(), Tags: tags, Type: "count"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Min() / int64(du), Tags: tags, Type: "min"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: t.Max() / int64(du), Tags: tags, Type: "max"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Mean() / du), Tags: tags, Type: "mean"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[2] / du), Tags: tags, Type: "95-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(ps[3] / du), Tags: tags, Type: "99-percentile"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate1()), Tags: tags, Type: "one-minute"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate5()), Tags: tags, Type: "five-minute"})
			sendBody.Add(&INFLUXMetric{Metric: name, Timestamp: now, Value: int64(t.Rate15()), Tags: tags, Type: "fifteen-minute"})
			metric.Clear()
		}
	})

	if sendBody.Len() > 0 {
		if len(c.Send) < c.MaxHttpOutstanding {
			c.Send <- sendBody
		} else {
			if !c.Quiet {
				fmt.Printf("Dropping flow: Q at %d\n", len(c.Send))
			}
		}
	}

	return nil
}

func (m *INFLUXMetric) ToWire() []byte {
	tags := make([]string, 0, (len(m.Tags))+2)
	tags = append(tags, m.Metric)
	tags = append(tags, "type="+influxEscape(m.Type))

	for k, v := range m.Tags {
		tags = append(tags, fmt.Sprintf("%s=%s", influxEscape(k), influxEscape(v)))
	}

	return []byte(fmt.Sprintf("%s value=%d %d\n", strings.Join(tags, ","), m.Value, m.Timestamp))
}

func (m *INFLUXMetricSet) Len() int {
	return len(m.metrics)
}

func (m *INFLUXMetricSet) ToWire() ([]byte, error) {
	var buf bytes.Buffer
	for _, met := range m.metrics {
		buf.Write(met.ToWire())
	}
	return buf.Bytes(), nil
}

// toWireChunks is much like ToWire, except it returns a slice of byte slices
// instead of just one byte slice. Assuming that maxChunkSize >
// len(InfluxMetric.ToWire()) for all metrics in the set, each byte slice
// returned by ToWireChunks is smaller than maxChunkSize. Any metric in the set
// that is larger than maxChunkSize will be in its own chunk.
func (m *INFLUXMetricSet) toWireChunks(maxChunkSize int) ([][]byte, error) {
	var buf bytes.Buffer
	chunks := make([][]byte, 0)

	saveChunk := func() {
		chunk := make([]byte, buf.Len())
		copy(chunk, buf.Bytes())
		buf.Reset()
		chunks = append(chunks, chunk)
	}

	for _, met := range m.metrics {
		metWire := met.ToWire()
		if buf.Len()+len(metWire) >= maxChunkSize {
			saveChunk()
		}
		buf.Write(metWire)
	}

	if buf.Len() > 0 {
		saveChunk()
	}
	return chunks, nil
}

func (m *INFLUXMetricSet) Add(met *INFLUXMetric) {
	if m.keepZero || met.Value > 0 {
		m.metrics = append(m.metrics, met)
	}
}

func NewINFLUXMetricSet(keepZero bool) *INFLUXMetricSet {
	return &INFLUXMetricSet{
		metrics:  make([]*INFLUXMetric, 0),
		keepZero: keepZero,
	}
}
