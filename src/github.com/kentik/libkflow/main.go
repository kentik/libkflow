package main

// #include "kflow.h"
import "C"
import (
	"fmt"
	"net/url"
	"reflect"
	"time"
	"unsafe"

	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
)

var sender *Sender
var errors chan error

//export kflowInit
func kflowInit(cfg *C.kflowConfig) C.int {
	errors = make(chan error, 100)

	url, err := url.Parse(C.GoString(cfg.URL))
	if err != nil {
		errors <- err
		return C.EKFLOWCONFIG
	}

	var (
		email   = C.GoString(cfg.API.email)
		token   = C.GoString(cfg.API.token)
		timeout = time.Duration(cfg.timeout) * time.Millisecond
		device  *api.Device
	)

	client := api.NewClient(email, token, timeout)

	switch url := C.GoString(cfg.API.URL); {
	case cfg.device_id > 0:
		device, err = client.GetDeviceByID(url, int(cfg.device_id))
	case cfg.hostname != nil:
		device, err = client.GetDeviceByName(url, C.GoString(cfg.hostname))
	default:
		err = fmt.Errorf("config: missing device selector")
	}

	if err != nil {
		errors <- err
		return C.EKFLOWCONFIG
	}

	interval := time.Duration(cfg.metrics.interval) * time.Second
	metrics := NewMetrics(device.ClientID())
	metrics.Start(C.GoString(cfg.metrics.URL), email, token, interval)

	agg, err := agg.NewAgg(time.Second, device.MaxFlowRate, &metrics.Metrics)
	if err != nil {
		errors <- err
		return C.EKFLOWCONFIG
	}

	sender = NewSender(url, timeout, int(cfg.verbose))
	sender.Errors = errors

	if err = sender.Start(agg, client, device, 2); err != nil {
		errors <- err
		sender = nil
		return C.EKFLOWCONFIG
	}

	return 0
}

//export kflowSend
func kflowSend(cflow *C.kflow) C.int {
	if sender == nil {
		return C.EKFLOWNOINIT
	}

	customs := *(*[]C.kflowCustom)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (uintptr)(unsafe.Pointer(cflow.customs)),
		Len:  int(cflow.numCustoms),
		Cap:  int(cflow.numCustoms),
	}))

	for i, c := range customs {
		name := C.GoString(c.name)
		id, ok := sender.Device.Customs[name]
		if !ok {
			return C.EKFLOWNOCUSTOM
		}
		customs[i].id = (C.uint64_t)(id)
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

func main() {
}
