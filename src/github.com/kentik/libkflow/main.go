package main

// #include "kflow.h"
import "C"
import (
	"net/url"
	"time"
)

var sender *Sender

//export kflowInit
func kflowInit(cfg *C.kflowConfig) C.int {
	url, err := url.Parse(C.GoString(cfg.URL))
	if err != nil {
		return C.EKFLOWCONFIG
	}

	timeout := time.Duration(cfg.timeout) * time.Millisecond

	s := NewSender(url, timeout, int(cfg.verbose))

	email := C.GoString(cfg.API.email)
	token := C.GoString(cfg.API.token)

	err = s.Validate(C.GoString(cfg.API.URL), email, token, int(cfg.device_id))
	if err != nil {
		return C.EKFLOWCONFIG
	}

	sender = s

	return 0
}

//export kflowSend
func kflowSend(cflow *C.kflow) C.int {
	if sender == nil {
		return C.EKFLOWNOINIT
	}

	msg, err := Pack((*Ckflow)(cflow))
	if err != nil {
		return C.EKFLOWNOMEM
	}

	if !sender.Send(msg) {
		return C.EKFLOWSEND
	}

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

func main() {
}
