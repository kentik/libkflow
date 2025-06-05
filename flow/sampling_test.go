package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeSampleRate(t *testing.T) {
	tests := []struct {
		name            string
		flows           []Flow
		resampleRateAdj float32
		expectedFlows   []Flow
	}{
		{
			name: "basic",
			flows: []Flow{
				{SampleRate: 0, SampleAdj: false},
				{SampleRate: 1, SampleAdj: false},
				{SampleRate: 2, SampleAdj: true},
			},
			resampleRateAdj: 0,
			expectedFlows: []Flow{
				{SampleRate: 0, SampleAdj: true},
				{SampleRate: 100, SampleAdj: true},
				{SampleRate: 200, SampleAdj: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NormalizeSampleRate(tt.flows, tt.resampleRateAdj)
			assert.Equal(t, tt.expectedFlows, tt.flows)
		})
	}
}
