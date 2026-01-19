package libkflow

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kentik/libkflow/api/test"
	"github.com/kentik/libkflow/flow"
)

func TestDirectFlowSender_SendFlows(t *testing.T) {
	client, server, device, err := test.NewClientServer()
	require.NoError(t, err)
	url := server.URL(test.FLOW)

	directSender := NewDirectSender(client, device, url)

	expected1 := flow.Flow{
		DeviceId:  uint32(device.ID),
		SrcAs:     rand.Uint32(),
		DstAs:     rand.Uint32(),
		SampleAdj: true,
	}
	expected2 := flow.Flow{
		DeviceId:  uint32(device.ID),
		SrcAs:     rand.Uint32(),
		DstAs:     rand.Uint32(),
		SampleAdj: true,
	}

	flows := []flow.Flow{expected1, expected2}
	n, err := directSender.SendFlows(flows)
	require.NoError(t, err)

	msgs, err := receive(server)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(flows), msgs.Len())
	assert.Greater(t, n, int64(0))
	assert.Equal(t, flowToCHF(expected1, t).String(), msgs.At(0).String())
	assert.Equal(t, flowToCHF(expected2, t).String(), msgs.At(1).String())
}
