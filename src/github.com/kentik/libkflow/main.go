package main

// #include "kflow.h"
import "C"
import (
	"fmt"
	"net/url"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"unsafe"

	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
)

var sender *Sender
var errors chan error

//export kflowInit
func kflowInit(cfg *C.kflowConfig, customs **C.kflowCustom, n *C.uint32_t) C.int {
	errors = make(chan error, 100)

	flowurl, err := url.Parse(C.GoString(cfg.URL))
	if err != nil {
		fail("invalid flow URL: %s", err)
		return C.EKFLOWCONFIG
	}

	var (
		email   = C.GoString(cfg.API.email)
		token   = C.GoString(cfg.API.token)
		timeout = time.Duration(cfg.timeout) * time.Millisecond
		proxy   *url.URL
		device  *api.Device
	)

	if cfg.proxy.URL != nil {
		proxy, err = url.Parse(C.GoString(cfg.proxy.URL))
		if err != nil {
			fail("invalid proxy URL: %s", err)
			return C.EKFLOWCONFIG
		}
	}

	client := api.NewClient(email, token, timeout, proxy)

	switch apiurl := C.GoString(cfg.API.URL); {
	case cfg.device_id > 0:
		device, err = client.GetDeviceByID(apiurl, int(cfg.device_id))
	case cfg.hostname != nil:
		device, err = client.GetDeviceByName(apiurl, C.GoString(cfg.hostname))
	default:
		err = fmt.Errorf("config: missing device selector")
	}

	if err != nil {
		fail("device lookup error: %s", err)
		return C.EKFLOWCONFIG
	}

	populateCustoms(device, customs, n)

	interval := time.Duration(cfg.metrics.interval) * time.Second
	metrics := NewMetrics(device.ClientID())
	metrics.Start(C.GoString(cfg.metrics.URL), email, token, interval, proxy)

	agg, err := agg.NewAgg(time.Second, device.MaxFlowRate, &metrics.Metrics)
	if err != nil {
		fail("agg setup error: %s", err)
		return C.EKFLOWCONFIG
	}

	sender = NewSender(flowurl, timeout, int(cfg.verbose))
	sender.Errors = errors

	if err = sender.Start(agg, client, device, 2); err != nil {
		fail("send startup error: %s", err)
		sender = nil
		return C.EKFLOWCONFIG
	}

	signal.Ignore(syscall.SIGPIPE)

	return 0
}

//export kflowSend
func kflowSend(cflow *C.kflow) C.int {
	if sender == nil {
		return C.EKFLOWNOINIT
	}

	kflow, err := Pack(sender.Segment(), (*Ckflow)(cflow))
	if err != nil {
		errors <- err
		return C.EKFLOWNOMEM
	}

	sender.Send(&kflow)

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
		case "uint32":
			vtype = C.KFLOWCUSTOMU32
		case "float32":
			vtype = C.KFLOWCUSTOMF32
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
