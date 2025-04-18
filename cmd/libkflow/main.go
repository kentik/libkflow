package main

// #include "../../kflow.h"
import "C"
import (
	"fmt"
	"net"
	"net/url"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	"github.com/kentik/libkflow"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/log"
	"github.com/kentik/libkflow/status"
)

var sender *libkflow.Sender
var server *status.Server
var errors chan error

//export kflowInit
func kflowInit(cfg *KflowConfig, dev *KflowDevice) C.int {
	errors = make(chan error, 100)

	if cfg == nil {
		return C.EKFLOWCONFIG
	}

	flowurl, err := url.Parse(C.GoString(cfg.URL))
	if err != nil {
		fail("invalid flow URL: %s", err)
		return C.EKFLOWCONFIG
	}

	apiurl, err := url.Parse(C.GoString(cfg.API.URL))
	if err != nil {
		fail("invalid API URL: %s", err)
		return C.EKFLOWCONFIG
	}

	metricsurl, err := url.Parse(C.GoString(cfg.metrics.URL))
	if err != nil {
		fail("invalid metrics URL: %s", err)
		return C.EKFLOWCONFIG
	}

	var (
		email   = C.GoString(cfg.API.email)
		token   = C.GoString(cfg.API.token)
		timeout = time.Duration(cfg.timeout) * time.Millisecond
		program = C.GoString(cfg.program)
		version = C.GoString(cfg.version)
		proxy   *url.URL
	)

	if program == "" || version == "" {
		return C.EKFLOWCONFIG
	}

	if cfg.proxy.URL != nil {
		proxy, err = url.Parse(C.GoString(cfg.proxy.URL))
		if err != nil {
			fail("invalid proxy URL: %s", err)
			return C.EKFLOWCONFIG
		}
	}

	config := libkflow.NewConfig(email, token, program, version)
	config.SetCapture(libkflow.Capture{
		Device:  C.GoString(cfg.capture.device),
		Snaplen: int32(cfg.capture.snaplen),
		Promisc: cfg.capture.promisc == 1,
	})
	config.SetProxy(proxy)
	config.SetSampleRate(int(cfg.sample_rate))
	config.SetTimeout(timeout)
	config.SetVerbose(int(cfg.verbose))
	config.OverrideURLs(apiurl, flowurl, metricsurl)
	config.SetLeveledLogger(log.Global())

	switch {
	case cfg.device_id > 0:
		did := int(cfg.device_id)
		sender, err = libkflow.NewSenderWithDeviceID(did, errors, config)
	case cfg.device_if != nil:
		dif := C.GoString(cfg.device_if)
		sender, err = libkflow.NewSenderWithDeviceIF(dif, errors, config)
	case cfg.device_ip != nil:
		dip := net.ParseIP(C.GoString(cfg.device_ip))
		sender, err = libkflow.NewSenderWithDeviceIP(dip, errors, config)
		if err == libkflow.ErrInvalidDevice {
			sender, err = tryCreateDevice(cfg, errors, config)
		}
	case cfg.device_name != nil:
		name := C.GoString(cfg.device_name)
		sender, err = libkflow.NewSenderWithDeviceName(name, errors, config)
		if err == libkflow.ErrInvalidDevice {
			sender, err = tryCreateDevice(cfg, errors, config)
		}
	default:
		fail("no device identifier supplied")
		return C.EKFLOWCONFIG
	}

	if err != nil {
		switch err {
		case libkflow.ErrInvalidAuth:
			return C.EKFLOWAUTH
		case libkflow.ErrInvalidDevice:
			return C.EKFLOWNODEVICE
		default:
			fail("library setup error: %s", err)
			return C.EKFLOWCONFIG
		}
	}

	if cfg.status.port != 0 {
		var host = C.GoString(cfg.status.host)
		var port = int(cfg.status.port)
		server = status.NewServer(host, port)
		go server.Start(sender.Metrics)
	}

	if cfg.dns.enable != 0 {
		url, err := url.Parse(C.GoString(cfg.dns.URL))
		if err != nil {
			fail("invalid DNS URL: %s", err)
			return C.EKFLOWCONFIG
		}

		interval := time.Duration(cfg.dns.interval)
		sender.StartDNS(url, interval*time.Second)
	}

	dev.id = C.uint64_t(sender.Device.ID)
	dev.name = C.CString(sender.Device.Name)
	dev.sample_rate = C.uint64_t(sender.Device.SampleRate)
	populateCustoms(sender.Device, &dev.customs, &dev.num_customs)

	signal.Ignore(syscall.SIGPIPE)

	return 0
}

//export kflowSend
func kflowSend(cflow *C.kflow) C.int {
	if sender == nil {
		return C.EKFLOWNOINIT
	}

	ckflow := (*flow.Ckflow)(unsafe.Pointer(cflow))
	flow := flow.New(ckflow)
	sender.Send(&flow)

	return 0
}

//export kflowStop
func kflowStop(msec C.int) C.int {
	if sender == nil {
		return C.EKFLOWNOINIT
	}

	wait := time.Duration(msec) * time.Millisecond
	if !sender.Stop(wait) {
		return C.EKFLOWTIMEOUT
	}
	return 0
}

//export kflowError
func kflowError() *C.char {
	select {
	case err := <-errors:
		return C.CString(err.Error())
	default:
		return nil
	}
}

//export kflowVersion
func kflowVersion() *C.char {
	return C.CString(libkflow.Version)
}

//export kflowSendDNS
func kflowSendDNS(q KflowDomainQuery, a *KflowDomainAnswer, n C.size_t) C.int {
	bytes := func(b C.kflowByteSlice) []byte {
		return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len))
	}

	question := api.DNSQuestion{
		Name: string(bytes(q.name)),
		Host: net.IP(bytes(q.host)),
	}

	answers := make([]api.DNSResourceRecord, n)
	for i, a := range (*[1 << 30]KflowDomainAnswer)(unsafe.Pointer(a))[:n:n] {
		answers[i] = api.DNSResourceRecord{
			IP:  net.IP(bytes(a.ip)),
			TTL: uint32(a.ttl),
		}
	}

	response := api.DNSResponse{
		Question: question,
		Answers:  answers,
	}

	sender.SendDNS(&response)

	return 0
}

//export kflowSendEncodedDNS
func kflowSendEncodedDNS(ptr *byte, len C.size_t) C.int {
	data := C.GoBytes(unsafe.Pointer(ptr), C.int(len))
	sender.SendEncodedDNS(data)
	return 0
}

func tryCreateDevice(cfg *KflowConfig, errors chan<- error, config *libkflow.Config) (*libkflow.Sender, error) {
	name := C.GoString(cfg.device_name)
	ip := net.ParseIP(C.GoString(cfg.device_ip))
	planID := int(cfg.device_plan)
	siteID := int(cfg.device_site)

	if ip == nil {
		ip = net.ParseIP(C.GoString(cfg.capture.ip))
	}

	if name == "" || ip == nil || planID == 0 {
		return nil, libkflow.ErrInvalidDevice
	}

	return libkflow.NewSenderWithNewDevice(&api.DeviceCreate{
		Name:        name,
		Type:        "host-nprobe-dns-www",
		Description: name,
		SampleRate:  1,
		BgpType:     "none",
		PlanID:      planID,
		SiteID:      siteID,
		Subtype:     "kprobe",
		IPs:         []net.IP{ip},
		CdnAttr:     "N",
	}, errors, config)
}

func populateCustoms(device *api.Device, ptr **C.kflowCustom, cnt *C.uint32_t) {
	if ptr == nil || cnt == nil {
		return
	}

	n := len(device.Customs)
	*ptr = (*C.kflowCustom)(C.calloc(C.size_t(n), C.sizeof_kflowCustom))
	*cnt = C.uint32_t(n)

	customs := *(*[]C.kflowCustom)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(*ptr)),
		Len:  int(n),
		Cap:  int(n),
	}))

	for i, c := range device.Customs {
		var vtype C.int
		switch c.Type {
		case "string":
			vtype = C.KFLOWCUSTOMSTR
		case "byte":
			vtype = C.KFLOWCUSTOMU8
		case "uint16":
			vtype = C.KFLOWCUSTOMU16
		case "uint32":
			vtype = C.KFLOWCUSTOMU32
		case "uint64":
			vtype = C.KFLOWCUSTOMU64
		case "int8":
			vtype = C.KFLOWCUSTOMI8
		case "int16":
			vtype = C.KFLOWCUSTOMI16
		case "int32":
			vtype = C.KFLOWCUSTOMI32
		case "int64":
			vtype = C.KFLOWCUSTOMI64
		case "float32":
			vtype = C.KFLOWCUSTOMF32
		case "float64":
			vtype = C.KFLOWCUSTOMF64
		case "addr":
			vtype = C.KFLOWCUSTOMADDR
		}

		customs[i] = C.kflowCustom{
			id:    C.uint64_t(c.ID),
			name:  C.CString(c.Name),
			vtype: vtype,
		}
	}
}

func fail(format string, args ...interface{}) {
	errors <- fmt.Errorf(format, args...)
}

func main() {
}

const (
	EKFLOWCONFIG   = C.EKFLOWCONFIG
	EKFLOWNOINIT   = C.EKFLOWNOINIT
	EKFLOWNOMEM    = C.EKFLOWNOMEM
	EKFLOWTIMEOUT  = C.EKFLOWTIMEOUT
	EKFLOWSEND     = C.EKFLOWSEND
	EKFLOWNOCUSTOM = C.EKFLOWNOCUSTOM
	EKFLOWAUTH     = C.EKFLOWAUTH
	EKFLOWNODEVICE = C.EKFLOWNODEVICE
)

type (
	KflowConfig       C.kflowConfig
	KflowCustom       C.kflowCustom
	KflowDevice       C.kflowDevice
	KflowDomainQuery  C.kflowDomainQuery
	KflowDomainAnswer C.kflowDomainAnswer
	KflowByteSlice    C.kflowByteSlice
)
