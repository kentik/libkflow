package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Header http.Header
	*http.Client
}

func NewClient(email, token string, timeout time.Duration, proxy *url.URL) *Client {
	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	header := http.Header{
		"X-CH-Auth-Email":     {email},
		"X-CH-Auth-API-Token": {token},
	}

	return &Client{
		Header: header,
		Client: client,
	}
}

func (c *Client) GetDeviceByID(url string, did int) (*Device, error) {
	return c.getdevice(fmt.Sprintf("%s/device/%d", url, did))
}

func (c *Client) GetDeviceByName(url string, name string) (*Device, error) {
	return c.getdevice(fmt.Sprintf("%s/device/%s", url, name))
}

func (c *Client) getdevice(url string) (*Device, error) {
	r, err := c.do("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, fmt.Errorf("api: HTTP status code %d", r.StatusCode)
	}

	dr := &DeviceResponse{}
	if err := json.NewDecoder(r.Body).Decode(dr); err != nil {
		return nil, err
	}

	return &dr.Device, nil
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
