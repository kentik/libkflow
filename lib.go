package libkflow

import (
	"errors"
	"fmt"
	"net"

	"github.com/kentik/libkflow/api"
)

var (
	ErrInvalidAuth   = errors.New("invalid API email/token")
	ErrInvalidConfig = errors.New("invalid config")
	ErrInvalidDevice = errors.New("invalid device")
)

// NewSenderWithDeviceID creates a new flow Sender given a device ID,
// error channel, and Config.
func NewSenderWithDeviceID(did int, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()

	d, err := lookupdev(client.GetDeviceByID(did))
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithDeviceIF creates a new flow Sender given a device interface name,
// error channel, and Config.
func NewSenderWithDeviceIF(dif string, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()

	d, err := lookupdev(client.GetDeviceByIF(dif))
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithDeviceIP creates a new flow Sender given a device IP address,
// error channel, and Config.
func NewSenderWithDeviceIP(dip net.IP, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()

	d, err := lookupdev(client.GetDeviceByIP(dip))
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithDeviceName creates a new flow Sender given a device name address,
// error channel, and Config.
func NewSenderWithDeviceName(name string, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()

	d, err := lookupdev(client.GetDeviceByName(name))
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithDeviceNameWithErrors creates a new flow Sender given a device name address and Config.
// The channel is closed after Sender.Stop is called and all flow has been dispatched.
func NewSenderWithDeviceNameWithErrors(name string, cfg *Config) (*Sender, <-chan error, error) {
	client := cfg.client()
	d, err := lookupdev(client.GetDeviceByName(name))
	if err != nil {
		return nil, nil, err
	}

	return cfg.startWithInternalErrors(client, d)
}

// NewSenderWithNewDevice creates a new device given device creation parameters,
// and then creates a new flow Sender with that device, the error channel, and
// the Config.
func NewSenderWithNewDevice(dev *api.DeviceCreate, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()

	d, err := client.CreateDevice(dev)
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithNewDeviceWithErrors creates a new device and returns a new flow Sender and an error channel which will
// report errors generated from the Sender. The channel is closed after Sender.Stop is called and all flow has been dispatched
func NewSenderWithNewDeviceWithErrors(dev *api.DeviceCreate, cfg *Config) (*Sender, <-chan error, error) {
	client := cfg.client()
	d, err := client.CreateDevice(dev)
	if err != nil {
		return nil, nil, err
	}

	return cfg.startWithInternalErrors(client, d)
}

// NewSenderWithNewSiteAndDevice creates a new device and site then returns a flow Sender for that newly created device
func NewSenderWithNewSiteAndDevice(siteAndDevice *api.SiteAndDeviceCreate, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()
	d, err := client.CreateDeviceAndSite(siteAndDevice)
	if err != nil {
		return nil, err
	}

	return cfg.start(client, d, errors)
}

// NewSenderWithNewSiteAndDeviceWithErrors is the same as NewSenderWithNewSiteAndDeviceWithErrors except rather than
// passing in a channel to receive errors, a channel is returned by the function. The channel is closed after Sender.Stop is called
// and all flow has been dispatched
func NewSenderWithNewSiteAndDeviceWithErrors(siteAndDevice *api.SiteAndDeviceCreate, cfg *Config) (*Sender, <-chan error, error) {
	client := cfg.client()
	d, err := client.CreateDeviceAndSite(siteAndDevice)
	if err != nil {
		return nil, nil, err
	}

	return cfg.startWithInternalErrors(client, d)
}

// NewSenderFromDevice returns a Sender for an existing Device
func NewSenderFromDevice(d *api.Device, errors chan<- error, cfg *Config) (*Sender, error) {
	client := cfg.client()
	return cfg.start(client, d, errors)
}

// NewSenderFromDeviceWithErrors returns a Sender and an error channel for an existing Device
func NewSenderFromDeviceWithErrors(d *api.Device, cfg *Config) (*Sender, <-chan error, error) {
	client := cfg.client()
	return cfg.startWithInternalErrors(client, d)
}

func lookupdev(dev *api.Device, err error) (*api.Device, error) {
	if err != nil {
		switch {
		case api.IsErrorWithStatusCode(err, 401):
			return nil, ErrInvalidAuth
		case api.IsErrorWithStatusCode(err, 404):
			return nil, ErrInvalidDevice
		default:
			return nil, fmt.Errorf("device lookup error: %s", err)
		}
	}
	return dev, nil
}
