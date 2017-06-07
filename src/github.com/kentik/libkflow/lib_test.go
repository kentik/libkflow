package libkflow_test

import (
	"net/url"
	"testing"

	"github.com/kentik/libkflow"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

func TestNewSenderWithDeviceID(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token)
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderWithDeviceID(dev.ID, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func TestNewSenderWithDeviceIP(t *testing.T) {
	dev, assert := setupLibTest(t)

	errors := make(chan error, 100)
	config := libkflow.NewConfig(email, token)
	config.OverrideURLs(apiurl, flowurl, metricsurl)

	s, err := libkflow.NewSenderWithDeviceIP(dev.IP, errors, config)

	assert.NotNil(s)
	assert.Nil(err)
}

func setupLibTest(t *testing.T) (*api.Device, *assert.Assertions) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	apiurl, _ = url.Parse(server.URL() + "/api/v5")
	flowurl, _ = url.Parse(server.URL() + "/chf")
	metricsurl, _ = url.Parse(server.URL() + "/tsdb")

	email = client.Header["X-CH-Auth-Email"][0]
	token = client.Header["X-CH-Auth-API-Token"][0]

	return device, assert
}

var (
	apiurl     *url.URL
	flowurl    *url.URL
	metricsurl *url.URL
	email      string
	token      string
)
