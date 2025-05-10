package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	ClientConfig
	deviceURL    string
	updateURL    string
	statusURL    string
	allDeviceURL string
	*http.Client
}

type ClientConfig struct {
	Email   string
	Token   string
	Timeout time.Duration
	Retries int
	API     *url.URL
	Proxy   *url.URL

	Logger interface{}
}

type ExportStatus struct {
	ID      int    `json:"export_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	EXPORT_STATUS_OK    = "OK"
	EXPORT_STATUS_ERROR = "ERROR"
	EXPORT_STATUS_START = "START"
	EXPORT_STATUS_HALT  = "HALT"
)

func NewClient(config ClientConfig) *Client {
	transport := &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = transport
	retryClient.RetryWaitMin = config.Timeout
	retryClient.RetryMax = config.Retries
	retryClient.Logger = config.Logger

	client := retryClient.StandardClient() // *http.Client

	if config.Proxy != nil {
		transport.Proxy = http.ProxyURL(config.Proxy)
	}

	return &Client{
		ClientConfig: config,
		deviceURL:    config.API.String() + "/device/%v",
		updateURL:    config.API.String() + "/company/%v/device/%v/tags/snmp",
		statusURL:    config.API.String() + "/cloudExport/status/%v",
		allDeviceURL: config.API.String() + "/devices",
		Client:       client,
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

func (c *Client) GetAllDevices() ([]*Device, error) {
	return c.getdevices(c.allDeviceURL)
}

func (c *Client) GetDevices(filter *DevicesFilter) ([]*Device, error) {
	type queryParams struct {
		key   string
		value string
	}
	qp := make([]queryParams, 0)
	if filter.FilterCloud {
		qp = append(qp, queryParams{"filterCloud", "true"})
	}
	if filter.CloudOnly {
		qp = append(qp, queryParams{"cloudOnly", "true"})
	}
	for _, v := range filter.Subtypes {
		qp = append(qp, queryParams{"subtypes[]", v})
	}
	for _, v := range filter.AugmentWith {
		qp = append(qp, queryParams{"augmentWith[]", v})
	}
	for _, v := range filter.Columns {
		qp = append(qp, queryParams{"columns[]", v})
	}

	if len(qp) > 0 {
		var buf strings.Builder
		for _, p := range qp {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(p.key) // don't escape to allow for array syntax (normally escaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(p.value))
		}

		return c.getdevices(c.allDeviceURL + "?" + buf.String())
	} else {
		return c.getdevices(c.allDeviceURL)
	}
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
	r, err := c.do("GET", url, "application/json", nil, false)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, c.error(r)
	}

	dw := &DeviceWrapper{}
	if err := json.NewDecoder(r.Body).Decode(dw); err != nil {
		return nil, err
	}

	return dw.Device, nil
}

func (c *Client) getdevices(url string) ([]*Device, error) {
	r, err := c.do("GET", url, "application/json", nil, false)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, c.error(r)
	}

	dw := &AllDeviceWrapper{}
	if err := json.NewDecoder(r.Body).Decode(dw); err != nil {
		return nil, err
	}

	return dw.Devices, nil
}

func (c *Client) CreateDeviceAndSite(siteDevCreate *SiteAndDeviceCreate) (*Device, error) {

	if !siteDevCreate.Device.AllowNoIP && len(siteDevCreate.Device.IPs) == 0 {
		return nil, fmt.Errorf("Missing IP for device")
	}

	createUrl := c.API.String() + "/deviceAndSite"

	// Remove chars which result in 500
	siteDevCreate.Device.NormalizeName()

	body, err := json.Marshal(siteDevCreate)
	if err != nil {
		return nil, err
	}

	r, err := c.do("POST", createUrl, "application/json", bytes.NewBuffer(body), false)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 && r.StatusCode != 201 {
		return nil, c.error(r)
	}

	dw := &DeviceWrapper{}
	if err := json.NewDecoder(r.Body).Decode(dw); err != nil {
		return nil, fmt.Errorf("response decoding error: %v", err)
	}

	// device structure returned from create call doesn't include necessary
	// fields like custom columns, so make another API call.

	return c.GetDeviceByID(dw.Device.ID)
}

func (c *Client) CreateDevice(create *DeviceCreate) (*Device, error) {
	url := fmt.Sprintf(c.deviceURL, "")

	if !create.AllowNoIP && len(create.IPs) == 0 {
		return nil, fmt.Errorf("Missing IP for device")
	}

	// Remove chars which result in 500
	create.NormalizeName()

	body, err := json.Marshal(map[string]*DeviceCreate{
		"device": create,
	})

	if err != nil {
		return nil, err
	}

	r, err := c.do("POST", url, "application/json", bytes.NewBuffer(body), false)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 201 {
		return nil, c.error(r)
	}

	dw := &DeviceWrapper{}
	if err := json.NewDecoder(r.Body).Decode(dw); err != nil {
		return nil, err
	}

	// device structure returned from create call doesn't include necessary
	// fields like custom columns, so make another API call.

	return c.GetDeviceByID(dw.Device.ID)
}

func (c *Client) GetInterfaces(did int) ([]Interface, error) {
	url := fmt.Sprintf(c.deviceURL+"/interfaces", did)

	r, err := c.do("GET", url, "application/json", nil, false)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, c.error(r)
	}

	interfaces := []Interface{}
	err = json.NewDecoder(r.Body).Decode(&interfaces)

	return interfaces, err
}

func (c *Client) UpdateInterfaces(dev *Device, nif *net.Interface) error {
	difs, err := c.GetInterfaces(dev.ID)
	if err != nil {
		return err
	}

	updates, err := GetInterfaceUpdates(nif)
	if err != nil {
		return err
	}

	if0 := InterfaceUpdate{
		Index: 0,
		Desc:  "kernel",
	}
	updates[if0.Desc] = if0

	for _, dif := range difs {
		name := dif.Desc
		if nif, ok := updates[name]; ok {
			if nif.Index == dif.Index &&
				nif.Desc == dif.Desc &&
				nif.Address == dif.Address &&
				nif.Netmask == dif.Netmask &&
				reflect.DeepEqual(nif.Addrs, dif.Addrs) {
				delete(updates, name)
			}
		}
	}

	if len(updates) == 0 {
		return nil
	}

	url := fmt.Sprintf(c.updateURL, dev.CompanyID, dev.ID)

	body, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	r, err := c.do("PUT", url, "application/json", bytes.NewBuffer(body), false)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	io.Copy(io.Discard, r.Body)

	if r.StatusCode != 200 {
		return &Error{StatusCode: r.StatusCode}
	}

	return nil
}

func (c *Client) UpdateInterfacesDirectly(dev *Device, updates map[string]InterfaceUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	url := fmt.Sprintf(c.updateURL, dev.CompanyID, dev.ID)

	body, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	r, err := c.do("PUT", url, "application/json", bytes.NewBuffer(body), false)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	io.Copy(io.Discard, r.Body)

	if r.StatusCode != 200 {
		return &Error{StatusCode: r.StatusCode}
	}

	return nil
}

func (c *Client) SendFlow(url string, buf *bytes.Buffer) error {
	r, err := c.do("POST", url, "application/binary", buf, true)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	_, _ = io.Copy(io.Discard, r.Body)

	if r.StatusCode != 200 {
		return fmt.Errorf("api: HTTP status code %d", r.StatusCode)
	}

	return nil
}

func (e *ExportStatus) Set(s string, m string) *ExportStatus {
	e.Status = s
	e.Message = m
	return e
}

// /internal/cloudExport/status/:id
// body expects `{status: string, message: string}`
func (c *Client) UpdateExportStatus(status *ExportStatus) error {
	url := fmt.Sprintf(c.statusURL, status.ID)

	body, err := json.Marshal(status)
	if err != nil {
		return err
	}

	r, err := c.do("PUT", url, "application/json", bytes.NewBuffer(body), false)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	io.Copy(io.Discard, r.Body)

	if r.StatusCode != 200 {
		return &Error{StatusCode: r.StatusCode}
	}

	return nil
}

func (c *Client) SendDNS(url string, buf *bytes.Buffer) error {
	r, err := c.do("POST", url, "application/chfdns", buf, false)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	io.Copy(io.Discard, r.Body)

	if r.StatusCode != 200 {
		return fmt.Errorf("api: HTTP status code %d", r.StatusCode)
	}

	return nil
}

func (c *Client) do(method, url, ctype string, body io.Reader, isGz bool) (*http.Response, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("X-CH-Auth-Email", c.Email)
	r.Header.Set("X-CH-Auth-API-Token", c.Token)
	r.Header.Set("Content-Type", ctype)
	if isGz {
		r.Header.Set("Content-Encoding", "gzip")
	}

	return c.Client.Do(r)
}

func (c *Client) error(r *http.Response) error {
	body := map[string]string{}
	json.NewDecoder(r.Body).Decode(&body)
	return &Error{
		StatusCode: r.StatusCode,
		Message:    body["error"],
	}
}
