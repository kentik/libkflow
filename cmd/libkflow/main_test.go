package main

import (
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var (
		n   int
		dev KflowDevice
	)
	cfg, apidev, assert := setupMainTest(t)

	// init with device ID
	n = int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	// init with device IP
	cfg.device_id = 0
	n = int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	// init with device name
	cfg.device_ip = nil
	n = int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	assert.Equal(apidev.ID, int(dev.id))
	assert.Equal(apidev.Name, cstr(dev.name))
	assert.Equal(apidev.SampleRate, int(dev.sample_rate))
}

func TestInitCreateDeviceWithName(t *testing.T) {
	var dev KflowDevice
	cfg, apidev, assert := setupMainTest(t)

	name := test.RandStr(8)
	ip := randip()
	devicename = append([]byte(name), 0)
	deviceip = append([]byte(ip.String()), 0)

	cfg.capture.ip = (*_Ctype_char)(unsafe.Pointer(&deviceip[0]))
	cfg.device_id = 0
	cfg.device_ip = nil
	cfg.device_name = (*_Ctype_char)(unsafe.Pointer(&devicename[0]))
	cfg.device_plan = _Ctype_int(rand.Uint32())

	n := int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	assert.Equal(name, cstr(dev.name))
	assert.NotEqual(apidev.Name, cstr(dev.name))
}

func TestInitCreateDeviceWithIP(t *testing.T) {
	var dev KflowDevice
	cfg, apidev, assert := setupMainTest(t)

	name := test.RandStr(8)
	ip := randip()

	deviceip = append([]byte(ip.String()), 0)
	devicename = append([]byte(name), 0)

	cfg.device_id = 0
	cfg.device_ip = (*_Ctype_char)(unsafe.Pointer(&deviceip[0]))
	cfg.device_name = (*_Ctype_char)(unsafe.Pointer(&devicename[0]))
	cfg.device_plan = _Ctype_int(rand.Uint32())

	n := int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	assert.Equal(name, cstr(dev.name))
	assert.NotEqual(apidev.Name, cstr(dev.name))
}

func TestInitInvalidConfig(t *testing.T) {
	var (
		n   int
		dev KflowDevice
	)
	assert := assert.New(t)

	// NULL config
	n = int(kflowInit(nil, &dev))
	assert.Equal(EKFLOWCONFIG, n)

	// NULL API URL
	cfg := KflowConfig{}
	n = int(kflowInit(&cfg, &dev))
	assert.Equal(EKFLOWCONFIG, n)
}

func TestInitMissingProgram(t *testing.T) {
	cfg, _, assert := setupMainTest(t)
	cfg.program = nil
	n := int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWCONFIG, n)
}

func TestInitMissingVersion(t *testing.T) {
	cfg, _, assert := setupMainTest(t)
	cfg.version = nil
	n := int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWCONFIG, n)
}

func TestInitInvalidAuth(t *testing.T) {
	cfg, _, assert := setupMainTest(t)
	cfg.API.email = nil
	n := int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWAUTH, n)
}

func TestInitInvalidDevice(t *testing.T) {
	var n int
	cfg, _, assert := setupMainTest(t)

	// invalid device ID
	cfg.device_id = cfg.device_id + 1
	n = int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWNODEVICE, n)

	// invalid device IP
	cfg.device_id = 0
	cfg.device_name = nil
	cfg.device_ip = (*_Ctype_char)(unsafe.Pointer(&deviceip[1]))
	n = int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWNODEVICE, n)
}

func TestInitInvalidCreate(t *testing.T) {
	var n int
	cfg, _, assert := setupMainTest(t)

	devicename = append([]byte(test.RandStr(8)), 0)
	deviceip = append([]byte(randip().String()), 0)

	cfg.device_id = 0
	cfg.device_ip = (*_Ctype_char)(unsafe.Pointer(&deviceip[0]))
	cfg.device_name = nil
	cfg.device_plan = _Ctype_int(rand.Uint32())

	// missing device name
	n = int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWNODEVICE, n)

	// missing device IP
	cfg.device_ip = nil
	cfg.device_name = (*_Ctype_char)(unsafe.Pointer(&devicename[0]))
	n = int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWNODEVICE, n)

	// missing device plan
	cfg.device_ip = (*_Ctype_char)(unsafe.Pointer(&deviceip[0]))
	cfg.device_plan = 0
	n = int(kflowInit(cfg, &KflowDevice{}))
	assert.Equal(EKFLOWNODEVICE, n)
}

func TestInitStatusServer(t *testing.T) {
	cfg, _, assert := setupMainTest(t)

	statushost := append([]byte("localhost"), 0)
	cfg.status.host = (*_Ctype_char)(unsafe.Pointer(&statushost[0]))
	cfg.status.port = 62000

	var dev KflowDevice
	n := int(kflowInit(cfg, &dev))
	assert.Equal(0, n)

	r, err := http.Get("http://localhost:62000/v1/status")
	for n := 0; n < 10 && err != nil; n++ {
		time.Sleep(100 * time.Millisecond)
		r, err = http.Get("http://localhost:62000/v1/status")
	}
	assert.NoError(err)
	assert.Equal(200, r.StatusCode)
}

func TestPopulateCustoms(t *testing.T) {
	assert := assert.New(t)

	device := api.Device{
		Customs: []api.Column{
			{1, "string", "string"},
			{2, "byte", "byte"},
			{3, "uint16", "uint16"},
			{4, "uint32", "uint32"},
			{5, "uint64", "uint64"},
			{6, "int8", "int8"},
			{7, "int16", "int16"},
			{8, "int32", "int32"},
			{9, "int64", "int64"},
			{10, "float32", "float32"},
			{11, "float64", "float64"},
			{12, "addr", "addr"},
		},
	}

	var ptr *_Ctype_struct___1
	var len _Ctype_uint32_t
	populateCustoms(&device, &ptr, &len)

	columns := *(*[]KflowCustom)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(ptr)),
		Len:  int(len),
		Cap:  int(len),
	}))

	assert.Equal(columns[0].vtype, _Ctype_int(1))
	assert.Equal(columns[1].vtype, _Ctype_int(2))
	assert.Equal(columns[2].vtype, _Ctype_int(3))
	assert.Equal(columns[3].vtype, _Ctype_int(4))
	assert.Equal(columns[4].vtype, _Ctype_int(5))
	assert.Equal(columns[5].vtype, _Ctype_int(6))
	assert.Equal(columns[6].vtype, _Ctype_int(7))
	assert.Equal(columns[7].vtype, _Ctype_int(8))
	assert.Equal(columns[8].vtype, _Ctype_int(9))
	assert.Equal(columns[9].vtype, _Ctype_int(10))
	assert.Equal(columns[10].vtype, _Ctype_int(11))
	assert.Equal(columns[11].vtype, _Ctype_int(12))
}

func setupMainTest(t *testing.T) (*KflowConfig, *api.Device, *assert.Assertions) {
	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	apiurl = append([]byte(server.URL(test.API).String()), 0)
	email = append([]byte(client.Email))
	token = append([]byte(client.Token))
	deviceip = append([]byte(device.IP.String()), 0)
	devicename = append([]byte(device.Name), 0)
	program = append([]byte("test"), 0)
	version = append([]byte("0.0.1"), 0)

	cfg := KflowConfig{
		API: _Ctype_struct___5{
			email: (*_Ctype_char)(unsafe.Pointer(&email[0])),
			token: (*_Ctype_char)(unsafe.Pointer(&token[0])),
			URL:   (*_Ctype_char)(unsafe.Pointer(&apiurl[0])),
		},
		device_id:   _Ctype_int(device.ID),
		device_ip:   (*_Ctype_char)(unsafe.Pointer(&deviceip[0])),
		device_name: (*_Ctype_char)(unsafe.Pointer(&devicename[0])),
		program:     (*_Ctype_char)(unsafe.Pointer(&program[0])),
		version:     (*_Ctype_char)(unsafe.Pointer(&version[0])),
	}

	return &cfg, device, assert
}

func cstr(c *_Ctype_char) string {
	str := *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(c)),
		Len:  1 << 31,
	}))
	idx := strings.IndexByte(str, 0)
	return string(str[:idx])
}

func randip() net.IP {
	return net.IPv4(
		byte(1+rand.Uint32()%255),
		byte(1+rand.Uint32()%255),
		byte(1+rand.Uint32()%255),
		byte(1+rand.Uint32()%255),
	)
}

var (
	apiurl     []byte
	email      []byte
	token      []byte
	deviceip   []byte
	devicename []byte
	program    []byte
	version    []byte
	statushost []byte
)
