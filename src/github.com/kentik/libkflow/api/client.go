package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Header    http.Header
	deviceURL string
	*http.Client
}

type ClientConfig struct {
	Email   string
	Token   string
	Timeout time.Duration
	API     *url.URL
	Proxy   *url.URL
}

func NewClient(config ClientConfig) *Client {
	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	if config.Proxy != nil {
		transport.Proxy = http.ProxyURL(config.Proxy)
	}

	header := http.Header{
		"X-CH-Auth-Email":     {config.Email},
		"X-CH-Auth-API-Token": {config.Token},
	}

	return &Client{
		Header:    header,
		deviceURL: config.API.String() + "/device/%v",
		Client:    client,
	}
}

func (c *Client) GetDeviceByID(did int) (*Device, error) {
	return c.getdevice(fmt.Sprintf(c.deviceURL, did))
}

func (c *Client) GetDeviceByName(name string) (*Device, error) {
	return c.getdevice(fmt.Sprintf(c.deviceURL, NormalizeName(name)))
}

func (c *Client) GetDeviceByIP(ip net.IP) (*Device, error) {
	return c.getdevice(fmt.Sprintf(c.deviceURL, ip))
}

func (c *Client) GetDeviceByIF(name string) (*Device, error) {
	nif, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := nif.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			dev, err := c.GetDeviceByIP(ip)
			if err == nil {
				return dev, err
			}
		}
	}

	return nil, &Error{StatusCode: 404}
}

func (c *Client) getdevice(url string) (*Device, error) {
	r, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, &Error{StatusCode: r.StatusCode}
	}

	dr := &DeviceResponse{}
	if err := json.NewDecoder(r.Body).Decode(dr); err != nil {
		return nil, err
	}

	return dr.Device, nil
}

func (c *Client) SendFlow(url string, buf *bytes.Buffer) error {
	r, err := c.do("POST", url, buf)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return fmt.Errorf("api: HTTP status code %d", r.StatusCode)
	}

	return nil
}

func (c *Client) do(method, url string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = c.Header
	return c.Client.Do(r)
}
