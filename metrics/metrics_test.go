package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsConfig(t *testing.T) {

	program := "test"
	version := "0.0.1"

	metrics := New(1, 2, program, version)

	metrics.Unregister()

	metrics.reg.Each(func(name string, _ interface{}) {
		assert.Contains(t, name, "ver="+program+"-"+version)
		assert.Contains(t, name, "ft="+program)
		assert.Contains(t, name, "dt=libkflow")
		assert.Contains(t, name, "level=primary")
		assert.Contains(t, name, "cid=1")
		assert.Contains(t, name, "did=2")
	})
}
