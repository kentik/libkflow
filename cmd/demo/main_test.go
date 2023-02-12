package main

import (
	"net"
	"sync/atomic"
	"testing"

	"github.com/kentik/libkflow"
	"github.com/stretchr/testify/assert"
)

type stubLogger struct {
	count uint32
}

func (s *stubLogger) Printf(string, ...interface{}) { atomic.AddUint32(&s.count, 1) }

type stubLeveledLogger struct {
	count uint32
}

func (s *stubLeveledLogger) Error(string, ...interface{}) { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Info(string, ...interface{})  { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Debug(string, ...interface{}) { atomic.AddUint32(&s.count, 1) }
func (s *stubLeveledLogger) Warn(string, ...interface{})  { atomic.AddUint32(&s.count, 1) }

func TestSettingLogger(t *testing.T) {
	// configuration here should not match an actual running server, but instead intentionally be down to verify logging
	var (
		email    = "test@example.com"
		token    = "token"
		deviceID = 1
		host     = net.ParseIP("127.0.0.1")
		port     = 8999
		program  = "demo"
		version  = "0.0.1"
	)
	errors := make(chan error, 100)

	config := libkflow.NewConfig(email, token, program, version)
	config.SetServer(host, port)
	config.SetRetries(1)
	config.SetVerbose(1)
	config.SetTimeout(0)

	l := stubLogger{}
	config.SetLogger(&l)

	// Performing the call, regardless of success, should result in logging internally
	_, _ = libkflow.NewSenderWithDeviceID(deviceID, errors, config)
	assert.True(t, l.count > 0)
}
func TestSettingLeveledLogger(t *testing.T) {
	// configuration here should not match an actual running server, but instead intentionally be down to verify logging
	var (
		email    = "test@example.com"
		token    = "token"
		deviceID = 1
		host     = net.ParseIP("127.0.0.1")
		port     = 8999
		program  = "demo"
		version  = "0.0.1"
	)
	errors := make(chan error, 100)

	config := libkflow.NewConfig(email, token, program, version)
	config.SetServer(host, port)
	config.SetRetries(1)
	config.SetVerbose(1)
	config.SetTimeout(0)

	l := stubLeveledLogger{}
	config.SetLeveledLogger(&l)

	// Performing the call, regardless of success, should result in logging internally
	_, _ = libkflow.NewSenderWithDeviceID(deviceID, errors, config)
	assert.True(t, l.count > 0)
}
