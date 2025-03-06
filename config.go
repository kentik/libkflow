package libkflow

import (
	"fmt"
	go_log "log"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	go_metrics "github.com/kentik/kit/go/legacy/go-metrics"
	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/log"
	"github.com/kentik/libkflow/metrics"
)

// Config describes the libkflow configuration.
type Config struct {
	email               string
	token               string
	capture             Capture
	proxy               *url.URL
	api                 *url.URL
	flow                *url.URL
	metrics             *url.URL
	sample              int
	timeout             time.Duration
	tlsHandshakeTimeout time.Duration
	retries             int
	logger              interface{}
	program             string
	version             string
	registry            go_metrics.Registry
	useInternalErrors   bool

	metricsPrefix   string
	metricsInterval time.Duration
}

// Capture describes the packet capture settings.
type Capture struct {
	Device  string
	Snaplen int32
	Promisc bool
}

// Logger interface allows to use other loggers than
// standard log.Logger.
type Logger interface {
	Printf(string, ...interface{})
}

// LeveledLogger interface implements the basic methods that a logger library needs
type LeveledLogger interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Warn(string, ...interface{})
}

// NewConfig returns a new Config given an API access email and token,
// and the name and version of the program using libkflow.
func NewConfig(email, token, program, version string) *Config {
	return defaultConfig(email, token, program, version)
}

// SetCapture sets the packet capture details.
func (c *Config) SetCapture(capture Capture) {
	c.capture = capture
}

// SetProxy sets the HTTP proxy used for making API requests, sending
// flow, and sending metrics.
func (c *Config) SetProxy(url *url.URL) {
	c.proxy = url
}

// SetServer changes the host and port used for API requests, flow,
// and metrics.
func (c *Config) SetServer(host net.IP, port int) {
	base := "http://" + net.JoinHostPort(host.String(), strconv.Itoa(port))
	var (
		api     = parseURL(base + "/api/internal")
		flow    = parseURL(base + "/chf")
		metrics = parseURL(base + "/tsdb")
	)
	c.OverrideURLs(api, flow, metrics)
}

// SetTimeout sets the HTTP request timeout.
func (c *Config) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// SetTLSHandshakeTimeout sets the TLSHandshakeTimeout on the http client's Transport
func (c *Config) SetTLSHandshakeTimeout(timeout time.Duration) {
	c.tlsHandshakeTimeout = timeout
}

// SetRetries sets the number of times to try HTTP requests.
func (c *Config) SetRetries(retries int) {
	c.retries = retries
}

// SetLogger sets the logger to use for the underlying HTTP requests.
func (c *Config) SetLogger(logger Logger) {
	c.logger = logger
}

// SetLeveledLogger sets the level based logger to use for the underlying HTTP requests.
func (c *Config) SetLeveledLogger(logger LeveledLogger) {
	c.logger = logger
}

// Set just the flow server
func (c *Config) SetFlow(server string) {
	c.flow = parseURL(server + "/chf")
}

// SetVerbose sets the verbosity level. Specifying a value greater than
// zero will cause verbose debug messages to be print to stderr.
func (c *Config) SetVerbose(level int) {
	log.SetVerbose(level)
}

// SetSampleRate sets the configured sample rate. If the sample rate
// is not set, and the rate configured in the device settings changes,
// then libkflow will abort the program with a call to exit().
func (c *Config) SetSampleRate(sample int) {
	c.sample = sample
}

// OverrideURLs changes the default endpoint URL for API requests,
// flow, and metrics.
func (c *Config) OverrideURLs(api, flow, metrics *url.URL) {
	c.api = api
	c.flow = flow
	c.metrics = metrics
}

func (c *Config) NewMetrics(dev *api.Device) *metrics.Metrics {
	return metrics.New(dev.CompanyID, dev.ID, c.program, c.version)
}

func (c *Config) GetClient() *api.Client {
	return c.client()
}

func (c *Config) SetMetricsPrefix(prefix string) {
	c.metricsPrefix = prefix
}

func (c *Config) SetMetricsInterval(dur time.Duration) {
	c.metricsInterval = dur
}

// WithRegistry allows setting a registry which will act as a shared registry between multiple Senders.
func (c *Config) WithRegistry(registry go_metrics.Registry) {
	c.registry = registry
}

func (c *Config) client() *api.Client {
	return api.NewClient(api.ClientConfig{
		Email:   c.email,
		Token:   c.token,
		Timeout: c.timeout,
		Retries: c.retries,
		API:     c.api,
		Proxy:   c.proxy,

		Logger: c.logger,
	})
}

func (c *Config) startWithInternalErrors(client *api.Client, dev *api.Device) (*Sender, <-chan error, error) {
	errChan := make(chan error)
	sender, err := c.start(client, dev, errChan)
	if err != nil {
		close(errChan)
		return nil, nil, err
	}

	sender.useInternalErrors = true

	return sender, errChan, nil
}

func (c *Config) start(client *api.Client, dev *api.Device, errors chan<- error) (*Sender, error) {
	if c.metricsInterval == 0 {
		c.metricsInterval = 60 * time.Second
	}

	var senderMetrics *metrics.Metrics
	if c.registry == nil {
		senderMetrics = c.NewMetrics(dev)
		senderMetrics.Start(c.metrics.String(), c.email, c.token, c.metricsPrefix, c.metricsInterval, c.proxy)
	} else {
		senderMetrics = metrics.NewWithRegistry(c.registry, dev.CompanyID, dev.ID, c.program, c.version)
	}

	agg, err := agg.NewAgg(time.Second, dev.MaxFlowRate, senderMetrics)
	if err != nil {
		return nil, fmt.Errorf("agg setup error: %s", err)
	}

	sender := newSender(c.flow, c.timeout)
	sender.Errors = errors
	sender.sample = c.sample
	sender.Metrics = senderMetrics

	if c.capture.Device != "" {
		nif, err := net.InterfaceByName(c.capture.Device)
		if err != nil {
			return nil, err
		}

		err = client.UpdateInterfaces(dev, nif)
		if err != nil {
			log.Debugf("error updating device interfaces: %s", err)
		}
	}

	if err = sender.start(agg, client, dev, 2); err != nil {
		return nil, fmt.Errorf("send startup error: %s", err)
	}

	return sender, nil
}

func defaultConfig(email, token, program, version string) *Config {
	return &Config{
		email:         email,
		token:         token,
		capture:       Capture{},
		proxy:         nil,
		api:           parseURL("https://api.kentik.com/api/internal"),
		flow:          parseURL("https://flow.kentik.com/chf"),
		metrics:       parseURL("https://flow.kentik.com/tsdb"),
		timeout:       10 * time.Second,
		retries:       0,
		logger:        go_log.New(os.Stderr, "", go_log.LstdFlags), // default behavior of underlying logger
		program:       program,
		version:       version,
		metricsPrefix: "chf",
	}
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic("invalid URL: " + s)
	}
	return u
}
