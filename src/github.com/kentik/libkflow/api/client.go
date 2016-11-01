package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Header http.Header
	*http.Client
}

func NewClient(email, token string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
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

func (c *Client) GetDevice(url string, did int) (*Device, error) {
	r, err := c.do("GET", fmt.Sprintf("%s/device/%d", url, did), nil)
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
