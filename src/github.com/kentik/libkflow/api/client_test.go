package api_test

import (
	"testing"

	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

func TestGetDeviceByID(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	device2, err := client.GetDeviceByID(device.ID)

	assert.NoError(err)
	assert.EqualValues(device, device2)
}

func TestGetDeviceByName(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	device2, err := client.GetDeviceByName(device.Name)

	assert.NoError(err)
	assert.EqualValues(device, device2)
}

func TestGetInvalidDevice(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	_, err = client.GetDeviceByName(device.Name + "-invalid")
	assert.Error(err)
	assert.Equal(&api.Error{StatusCode: 404}, err)
}
