package api_test

import (
	"math/rand"
	"net"
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

func TestGetDeviceByIP(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	device2, err := client.GetDeviceByIP(device.IP)

	assert.NoError(err)
	assert.EqualValues(device, device2)
}

func TestGetDeviceByIF(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	ifs, err := net.Interfaces()
	assert.NoError(err)

	device2, err := client.GetDeviceByIF(ifs[0].Name)

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
	assert.True(api.IsErrorWithStatusCode(err, 404))

	_, err = client.GetDeviceByIF("invalid")
	assert.Error(err)

	_, err = client.GetDeviceByIP(net.ParseIP("0.0.0.0"))
	assert.Error(err)
	assert.True(api.IsErrorWithStatusCode(err, 404))
}

func TestCreateDevice(t *testing.T) {
	client, _, _, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	create := &api.DeviceCreate{
		Name:        test.RandStr(8),
		Type:        test.RandStr(8),
		Description: test.RandStr(8),
		SampleRate:  int(rand.Uint32()),
		BgpType:     test.RandStr(4),
		PlanID:      int(rand.Uint32()),
		IPs:         []net.IP{net.ParseIP("127.0.0.1")},
		CdnAttr:     test.RandStr(1),
	}

	device, err := client.CreateDevice(create)
	assert.NoError(err)

	assert.EqualValues(create.Name, device.Name)
	assert.EqualValues(create.Type, device.Type)
	assert.EqualValues(create.Description, device.Description)
	assert.EqualValues(create.IPs[0], device.IP)
	assert.EqualValues(create.SampleRate, device.SampleRate)
	assert.EqualValues(create.BgpType, device.BgpType)
	assert.EqualValues(create.PlanID, int(device.Plan.ID))
	assert.EqualValues(create.CdnAttr, device.CdnAttr)
}

func TestCreateDeviceFailsWithoutIPFlag(t *testing.T) {
	client, _, _, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	create := &api.DeviceCreate{
		Name:        test.RandStr(8),
		Type:        test.RandStr(8),
		Description: test.RandStr(8),
		SampleRate:  int(rand.Uint32()),
		BgpType:     test.RandStr(4),
		PlanID:      int(rand.Uint32()),
		IPs:         []net.IP{},
		CdnAttr:     test.RandStr(1),
	}

	_, err = client.CreateDevice(create)
	assert.Error(err)
}

func TestCreateDeviceWithoutIP(t *testing.T) {
	client, _, _, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	create := &api.DeviceCreate{
		Name:        test.RandStr(8),
		Type:        test.RandStr(8),
		Description: test.RandStr(8),
		SampleRate:  int(rand.Uint32()),
		BgpType:     test.RandStr(4),
		PlanID:      int(rand.Uint32()),
		IPs:         []net.IP{},
		AllowNoIP:   true,
		CdnAttr:     test.RandStr(1),
	}

	device, err := client.CreateDevice(create)
	assert.NoError(err)

	assert.EqualValues(create.Name, device.Name)
	assert.EqualValues(create.Type, device.Type)
	assert.EqualValues(create.Description, device.Description)
	assert.EqualValues(net.IP(nil), device.IP)
	assert.EqualValues(create.SampleRate, device.SampleRate)
	assert.EqualValues(create.BgpType, device.BgpType)
	assert.EqualValues(create.PlanID, int(device.Plan.ID))
	assert.EqualValues(create.CdnAttr, device.CdnAttr)
}

func TestCreateDeviceAndSiteWithoutIP(t *testing.T) {
	client, _, _, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	create := &api.DeviceCreate{
		Name:        test.RandStr(8),
		Type:        test.RandStr(8),
		Description: test.RandStr(8),
		SampleRate:  int(rand.Uint32()),
		BgpType:     test.RandStr(4),
		PlanID:      int(rand.Uint32()),
		IPs:         []net.IP{},
		AllowNoIP:   true,
		CdnAttr:     test.RandStr(1),
	}

	siteCreate := &api.SiteCreate{
		Title:   "Hawaii Offsite",
		City:    "Honolulu",
		Region:  "Pacific",
		Country: "USA",
	}

	siteAndDeviceCreate := &api.SiteAndDeviceCreate{
		Site:   siteCreate,
		Device: create,
	}

	device, err := client.CreateDeviceAndSite(siteAndDeviceCreate)
	assert.NoError(err)

	assert.EqualValues(siteAndDeviceCreate.Device.Name, device.Name)
	assert.EqualValues(siteAndDeviceCreate.Device.Type, device.Type)
	assert.EqualValues(siteAndDeviceCreate.Device.Description, device.Description)
	assert.EqualValues(net.IP(nil), device.IP)
	assert.EqualValues(siteAndDeviceCreate.Device.SampleRate, device.SampleRate)
	assert.EqualValues(siteAndDeviceCreate.Device.BgpType, device.BgpType)
	assert.EqualValues(siteAndDeviceCreate.Device.PlanID, int(device.Plan.ID))
	assert.EqualValues(siteAndDeviceCreate.Device.CdnAttr, device.CdnAttr)
}

func TestGetAllDevices(t *testing.T) {
	client, _, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)

	devices, err := client.GetAllDevices()

	assert.NoError(err)
	assert.EqualValues(device, devices[0])
}
