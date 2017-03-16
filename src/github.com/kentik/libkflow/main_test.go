package main

import (
	"testing"
	"unsafe"

	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var n int
	cfg, assert := setupMainTest(t)

	// init with device ID
	n = int(kflowInit(cfg, nil, nil))
	assert.Equal(0, n)

	// init with hostname
	cfg.device_id = 0
	n = int(kflowInit(cfg, nil, nil))
	assert.Equal(0, n)
}

func TestInitInvalidConfig(t *testing.T) {
	var n int
	assert := assert.New(t)

	// NULL config
	n = int(kflowInit(nil, nil, nil))
	assert.Equal(EKFLOWCONFIG, n)

	// NULL API URL
	cfg := _Ctype_struct___3{}
	n = int(kflowInit(&cfg, nil, nil))
	assert.Equal(EKFLOWCONFIG, n)
}

func TestInitInvalidAuth(t *testing.T) {
	cfg, assert := setupMainTest(t)
	cfg.API.email = nil
	n := int(kflowInit(cfg, nil, nil))
	assert.Equal(EKFLOWAUTH, n)
}

func TestInitInvalidDevice(t *testing.T) {
	var n int
	cfg, assert := setupMainTest(t)

	// invalid device ID
	cfg.device_id = cfg.device_id + 1
	n = int(kflowInit(cfg, nil, nil))
	assert.Equal(EKFLOWNODEVICE, n)

	// invalid hostname
	cfg.device_id = 0
	cfg.hostname = (*_Ctype_char)(unsafe.Pointer(&hostname[1]))
	n = int(kflowInit(cfg, nil, nil))
	assert.Equal(EKFLOWNODEVICE, n)
}

func setupMainTest(t *testing.T) (*_Ctype_struct___3, *assert.Assertions) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	apiurl = append([]byte(server.URL()+"/api/v5"), 0)
	email = append([]byte(client.Header["X-CH-Auth-Email"][0]), 0)
	token = append([]byte(client.Header["X-CH-Auth-API-Token"][0]), 0)
	hostname = append([]byte(device.Name), 0)

	cfg := _Ctype_struct___3{
		API: _Ctype_struct___4{
			email: (*_Ctype_char)(unsafe.Pointer(&email[0])),
			token: (*_Ctype_char)(unsafe.Pointer(&token[0])),
			URL:   (*_Ctype_char)(unsafe.Pointer(&apiurl[0])),
		},
		device_id: _Ctype_int(device.ID),
		hostname:  (*_Ctype_char)(unsafe.Pointer(&hostname[0])),
	}

	return &cfg, assert
}

var (
	apiurl   []byte
	email    []byte
	token    []byte
	hostname []byte
)
