package cmetrics

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/kentik/kit/go/legacy/common/logging"
	"github.com/kentik/kit/go/legacy/go-metrics"
)

const (
	MAX_SEND_TRIES          = 2
	CLIENT_RESPONSE_TIMEOUT = 5 * time.Second
	CLIENT_KEEP_ALIVE       = 60 * time.Second
	CLIENT_TLS_TIMEOUT      = 5 * time.Second
	ContentType             = "application/json"
	API_EMAIL_HEADER        = "X-CH-Auth-Email"
	API_PASSWORD_HEADER     = "X-CH-Auth-API-Token"
	maxHttpTsdbBodyLen      = 1000 * 1000 // 1MB
)

type TSDBMetric struct {
	Metric    string            `json:"metric"`
	Timestamp int64             `json:"timestamp"`
	Value     int64             `json:"value"`
	Tags      map[string]string `json:"tags"`
}

// OpenHTTPTSDBConfig provides a container with configuration parameters for
// the OpenTSDB exporter
type OpenHTTPTSDBConfig struct {
	Addr               string            // Network address to connect to
	Registry           metrics.Registry  // Registry to be exported
	FlushInterval      time.Duration     // Flush interval
	DurationUnit       time.Duration     // Time conversion unit for durations
	Prefix             string            // Prefix to be prepended to metric names
	Debug              bool              // write to stdout for debug
	Tags               map[string]string // add these tags to each metric writen
	Send               chan []byte       // manage # of outstanding http requests here.
	MaxHttpOutstanding int
	ProxyUrl           string
	Extra              map[string]string
	ApiEmail           *string
	ApiPassword        *string
	l                  logging.Logger
	logPrefix          string
}

// OpenHTTPTSDBWithConfig is a blocking exporter function just like OpenHTTPTSDB,
// but it takes a OpenHTTPTSDBConfig instead.
func OpenHTTPTSDBWithConfig(c OpenHTTPTSDBConfig) {
	go c.runSend()

	for _ = range time.Tick(c.FlushInterval) {
		if err := openHTTPTSDB(&c); nil != err {
			c.l.Errorf(c.logPrefix, "Error opening HTTP TSDB: %s", err)
		}
	}
}

func addTypeTag(in map[string]string, mtype string) map[string]string {
	out := make(map[string]string)

	// Copy these over as a base.
	for k, v := range in {
		out[k] = v
	}

	// Add in type, and send
	out["type"] = mtype
	return out
}

func (c *OpenHTTPTSDBConfig) runSend() {
	tr := &http.Transport{
		DisableCompression: false,
		DisableKeepAlives:  false,
		Dial: (&net.Dialer{
			Timeout:   CLIENT_RESPONSE_TIMEOUT,
			KeepAlive: CLIENT_KEEP_ALIVE,
		}).Dial,
		TLSHandshakeTimeout: CLIENT_TLS_TIMEOUT,
	}

	// Add a proxy if needed.
	if c.ProxyUrl != "" {
		proxyUrl, err := url.Parse(c.ProxyUrl)
		if err != nil {
			c.l.Errorf(c.logPrefix, "Error setting proxy: %s", err)
		} else {
			tr.Proxy = http.ProxyURL(proxyUrl)
			c.l.Infof(c.logPrefix, "Set outbound proxy: %s", c.ProxyUrl)
		}
	}

	client := &http.Client{Transport: tr, Timeout: CLIENT_RESPONSE_TIMEOUT}

	for r := range c.Send {
		if c.Debug {
			c.l.Infof(c.logPrefix, "Metrics %s", string(r))
		} else {
			for i := 0; i < MAX_SEND_TRIES; i++ {
				req, err := http.NewRequest("POST", c.Addr, bytes.NewBuffer(r))
				if err != nil {
					c.l.Errorf(c.logPrefix, "Error Creating Request: %s", err)
					continue
				}

				req.Header.Add("Content-Type", ContentType)

				if c.ApiEmail != nil && c.ApiPassword != nil {
					req.Header.Add(API_EMAIL_HEADER, *c.ApiEmail)
					req.Header.Add(API_PASSWORD_HEADER, *c.ApiPassword)
				}

				resp, err := client.Do(req)
				if err != nil {
					if i > 0 {
						c.l.Errorf(c.logPrefix, "Error Posting to %s %s", c.Addr, err)
					} else {
						c.l.Warnf(c.logPrefix, "Retry Posting to %s %s", c.Addr, err)
					}
					client = &http.Client{Transport: tr, Timeout: CLIENT_RESPONSE_TIMEOUT}
				} else {
					// Fire and forget
					io.Copy(ioutil.Discard, resp.Body)
					resp.Body.Close()
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
func openHTTPTSDB(c *OpenHTTPTSDBConfig) error {
	shortHostnameBase := GetShortHostname()
	now := time.Now().Unix()
	sendBody := make([]TSDBMetric, 0)
	du := float64(c.DurationUnit)

	c.Registry.Each(func(baseName string, i interface{}) {
		name, tags := ExpandTags(baseName, c.Prefix, shortHostnameBase, c.Tags, c.Extra)

		switch metric := i.(type) {
		case metrics.Counter:
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: metric.Count(), Tags: addTypeTag(tags, "count")})
		case metrics.Gauge:
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: metric.Value(), Tags: addTypeTag(tags, "value")})
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: h.Count(), Tags: addTypeTag(tags, "count")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: h.Min(), Tags: addTypeTag(tags, "min")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: h.Max(), Tags: addTypeTag(tags, "max")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(h.Mean()), Tags: addTypeTag(tags, "mean")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(ps[2]), Tags: addTypeTag(tags, "95-percentile")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(ps[3]), Tags: addTypeTag(tags, "99-percentile")})
			metric.Clear()
		case metrics.Meter:
			m := metric.Snapshot()
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: m.Count(), Tags: addTypeTag(tags, "count")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(m.Rate1()), Tags: addTypeTag(tags, "one-minute")})
		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: t.Count(), Tags: addTypeTag(tags, "count")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: t.Min() / int64(du), Tags: addTypeTag(tags, "min")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: t.Max() / int64(du), Tags: addTypeTag(tags, "max")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(t.Mean() / du), Tags: addTypeTag(tags, "mean")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(ps[2] / du), Tags: addTypeTag(tags, "95-percentile")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(ps[3] / du), Tags: addTypeTag(tags, "99-percentile")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(t.Rate1()), Tags: addTypeTag(tags, "one-minute")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(t.Rate5()), Tags: addTypeTag(tags, "five-minute")})
			sendBody = append(sendBody, TSDBMetric{Metric: name, Timestamp: now, Value: int64(t.Rate15()), Tags: addTypeTag(tags, "fifteen-minute")})
			metric.Clear()
		}
	})

	sendBodyPruned := make([]TSDBMetric, 0)
	for _, m := range sendBody {
		if m.Value > 0 {
			sendBodyPruned = append(sendBodyPruned, m)
		}
	}

	if len(sendBodyPruned) == 0 {
		return nil
	}

	chunks, err := toChunks(sendBodyPruned, maxHttpTsdbBodyLen)
	if err != nil {
		c.l.Errorf(c.logPrefix, "Error encoding json: %s", err)
		return nil
	}

	for _, chunk := range chunks {
		if len(c.Send) < c.MaxHttpOutstanding {
			c.Send <- chunk
		} else {
			c.l.Errorf(c.logPrefix, "Dropping flow: Q at %d", len(c.Send))
		}
	}

	return nil
}

// toChunks encodes the provided metrics to JSON byte slices. Assuming that each metric's JSON is less than maxChunkSize,
// each slice's size will be smaller than maxChunkSize. If a metric's JSON is larger than maxChunkSize, it will be returned
// in its own chunk
func toChunks(metrics []TSDBMetric, maxChunkSize int) ([][]byte, error) {
	ret := make([][]byte, 0)

	for i := 0; i < len(metrics); {
		chunk, numWritten, err := createChunk(metrics[i:], maxChunkSize)
		if err != nil {
			return nil, err
		}

		ret = append(ret, chunk)

		i += numWritten
	}

	return ret, nil
}

// createChunk encodes one or more TSDBMetric objects from metrics as a JSON array and returns the resulting JSON and number of objects written.
// The returned byte slice will be smaller than maxSize except when the first metric's JSON is larger than maxSize,
// in which case that metric will be the only element in the returned JSON array.
func createChunk(metrics []TSDBMetric, maxSize int) ([]byte, int, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	numWritten := 0
	for _, metric := range metrics {
		metricBytes, err := json.Marshal(metric)
		if err != nil {
			return nil, 0, err
		}

		// the number of bytes needed to finish encoding the json array if we were to write the current
		// metric to the buffer
		var bytesToClose int
		if numWritten > 0 {
			bytesToClose = 2 // a comma and a closing bracket
		} else {
			bytesToClose = 1 // only a closing bracket
		}

		// check if we've written at least one object and if writing the current metric to the buffer
		// and finalizing the json encoding will bring the buffer over the maxSize
		if numWritten > 0 && buffer.Len()+len(metricBytes)+bytesToClose > maxSize {
			break
		}

		if numWritten > 0 {
			buffer.WriteString(",")
		}

		buffer.Write(metricBytes)
		numWritten++
	}

	buffer.WriteString("]")
	return buffer.Bytes(), numWritten, nil
}
