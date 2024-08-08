package libkflow

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/kentik/libkflow/agg"
	"github.com/kentik/libkflow/api"
	"github.com/kentik/libkflow/api/test"
	"github.com/kentik/libkflow/chf"
	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/metrics"
	"github.com/stretchr/testify/assert"
	capnp "zombiezen.com/go/capnproto2"
)

func TestSender(t *testing.T) {
	sender, server, assert := setup(t)

	expected := flow.Flow{
		DeviceId:  uint32(sender.Device.ID),
		SrcAs:     rand.Uint32(),
		DstAs:     rand.Uint32(),
		SampleAdj: true,
	}

	sender.Send(&expected)

	msgs, err := receive(server)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(sender.Stop(100 * time.Millisecond))
	assert.Greater(sender.Metrics.BytesSent.Count(), int64(0))
	assert.Equal(flowToCHF(expected, t).String(), msgs.At(0).String())
}

func TestSenderStop(t *testing.T) {
	sender, _, assert := setup(t)
	stopped := sender.Stop(100 * time.Millisecond)
	assert.True(stopped)
}

func TestSendDNS(t *testing.T) {
	sender, server, assert := setup(t)

	url := server.URL(test.DNS)
	sender.StartDNS(url, 1*time.Millisecond)

	expected := &api.DNSResponse{
		Question: api.DNSQuestion{
			Name: "foo.com",
			Host: net.ParseIP("127.0.0.1"),
		},
		Answers: []api.DNSResourceRecord{
			{
				Name:  "",
				CNAME: "",
				IP:    net.ParseIP("10.0.0.1"),
				TTL:   16,
			},
		},
	}

	sender.SendDNS(expected)

	select {
	case res := <-server.Dns():
		assert.Equal(expected, res)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("failed to receive DNS data")
	}
}

func BenchmarkSenderSend(b *testing.B) {
	sender, _, _ := setup(b)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		sender.Send(&flow.Flow{
			SrcAs: uint32(b.N),
			DstAs: uint32(b.N),
		})
	}
}

func setup(t testing.TB) (*Sender, *test.Server, *assert.Assertions) {
	metrics := metrics.New(100, 200, "send_test", "1.0.0")
	agg, err := agg.NewAgg(10*time.Millisecond, 100, metrics)
	if err != nil {
		t.Fatal(err)
	}

	client, server, device, err := test.NewClientServer()
	if err != nil {
		t.Fatal(err)
	}

	server.Log.SetOutput(io.Discard)

	url := server.URL(test.FLOW)
	sender := newSender(url, 1*time.Second)
	sender.Metrics = metrics
	sender.start(agg, client, device, 1)

	return sender, server, assert.New(t)
}

func receive(s *test.Server) (*chf.CHF_List, error) {
	interval := 100 * time.Millisecond
	select {
	case flow := <-s.Flows():
		msgs, err := flow.Msgs()
		return &msgs, err
	case <-time.After(interval):
		return nil, fmt.Errorf("failed to receive flow within %s", interval)
	}
}

func flowToCHF(flow flow.Flow, t testing.TB) chf.CHF {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	kflow, err := chf.NewCHF(seg)
	if err != nil {
		t.Fatal(err)
	}

	list, err := chf.NewCustom_List(seg, int32(len(flow.Customs)))
	if err != nil {
		t.Fatal(err)
	}

	flow.FillCHF(kflow, list)

	return kflow
}

func TestCompress(t *testing.T) {
	sender, server, assert := setup(t)

	expected := []flow.Flow{
		flow.Flow{
			DeviceId:  uint32(sender.Device.ID),
			SrcAs:     uint32(1),
			DstAs:     rand.Uint32(),
			SampleAdj: true,
		},
		flow.Flow{
			DeviceId:  uint32(sender.Device.ID),
			SrcAs:     uint32(2),
			DstAs:     rand.Uint32(),
			SampleAdj: true,
		},
		flow.Flow{
			DeviceId:  uint32(sender.Device.ID),
			SrcAs:     uint32(3),
			DstAs:     rand.Uint32(),
			SampleAdj: true,
		},
	}

	// Send them all at once
	for _, e := range expected {
		sender.Send(&e)
	}

	// Itterate through, looking at each one recieved.
	for i, e := range expected {
		msgs, err := receive(server)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(flowToCHF(e, t).String(), msgs.At(0).String(), "%d", i)
	}
}
