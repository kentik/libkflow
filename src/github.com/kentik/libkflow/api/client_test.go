package api_test

import (
	"testing"

	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

func TestGetDeviceByID(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	device2, err := client.GetDeviceByID(server.URL()+"/api/v5", device.ID)

	assert.NoError(err)
	assert.EqualValues(device, device2)
}

func TestGetDeviceByName(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	device2, err := client.GetDeviceByName(server.URL()+"/api/v5", device.Name)

	assert.NoError(err)
	assert.EqualValues(device, device2)
}
