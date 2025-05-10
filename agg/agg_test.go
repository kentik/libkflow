package agg

import (
	"bytes"
	"testing"
	"time"

	capnp "zombiezen.com/go/capnproto2"

	"github.com/stretchr/testify/assert"

	"github.com/kentik/libkflow/chf"
	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/metrics"
)

func TestAggSimple(t *testing.T) {
	var (
		interval  = 100 * time.Millisecond
		fps       = 100
		count     = fps / int(time.Second/interval)
		s, assert = setup(t, interval, fps)
	)

	flows := s.send(count, func(i int, flow *flow.Flow) {
		flow.SrcEthMac = uint64(i)
	})

	assert.Equal(count, len(flows))

	for i, flow := range flows {
		assert.EqualValues(i, flow.SrcEthMac())
	}

	assert.EqualValues(count, s.metrics.TotalFlowsIn.Count())
	assert.EqualValues(count, s.metrics.TotalFlowsOut.Count())
	assert.EqualValues(0, s.metrics.RateLimitDrops.Count())
}

func TestAggDrop(t *testing.T) {
	var (
		interval  = 100 * time.Millisecond
		fps       = 100
		buffer    = fps * MaxFlowBuffer
		expect    = fps / int(time.Second/interval)
		dropped   = 200
		count     = buffer/int(time.Second/interval) + dropped
		s, assert = setup(t, interval, fps)
	)

	flows := s.send(count, func(i int, flow *flow.Flow) {})

	assert.Equal(expect, len(flows))

	assert.EqualValues(count, s.metrics.TotalFlowsIn.Count())
	assert.EqualValues(expect, s.metrics.TotalFlowsOut.Count())
	assert.EqualValues(dropped, s.metrics.RateLimitDrops.Count())
}

type testState struct {
	interval time.Duration
	fps      int
	agg      *Agg
	metrics  *metrics.Metrics
	*testing.T
}

func setup(t *testing.T, interval time.Duration, fps int) (*testState, *assert.Assertions) {
	metrics := &metrics.Metrics{
		TotalFlowsIn:   metrics.NewMeter(),
		TotalFlowsOut:  metrics.NewMeter(),
		RateLimitDrops: metrics.NewMeter(),
		BytesSent:      metrics.NewMeter(),
	}

	agg, err := NewAgg(interval, fps, metrics)
	if err != nil {
		t.Fatal(err)
	}

	return &testState{
		interval: interval,
		fps:      fps,
		agg:      agg,
		metrics:  metrics,
		T:        t,
	}, assert.New(t)
}

func (s *testState) send(n int, g func(int, *flow.Flow)) []chf.CHF {
	for i := 0; i < n; i++ {
		f := flow.Flow{}
		g(i, &f)
		s.agg.Add(&f)
	}
	return s.receive()
}

func (s *testState) receive() []chf.CHF {
	interval := s.interval * 2

	buf := &bytes.Buffer{}

	select {
	case msg := <-s.agg.Output():
		err := capnp.NewPackedEncoder(buf).Encode(msg)
		if err != nil {
			s.Fatal(err)
		}
	case <-time.After(interval):
		s.Fatalf("failed to receive flow within %s", interval)
	}

	msg, err := capnp.NewPackedDecoder(buf).Decode()
	if err != nil {
		s.Fatal(err)
	}

	root, err := chf.ReadRootPackedCHF(msg)
	if err != nil {
		s.Fatal(err)
	}

	msgs, err := root.Msgs()
	if err != nil {
		s.Fatal(err)
	}

	flows := make([]chf.CHF, msgs.Len())
	for i := range flows {
		flows[i] = msgs.At(i)
	}

	return flows
}

func BenchmarkEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		flows := make([]flow.Flow, 100)

		for i := 0; i < len(flows); i++ {
			flows[i].DstAs = uint32(n)
			flows[i].SrcAs = uint32(n)
		}

		_, err := encode(flows)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func encode(flows []flow.Flow) (*capnp.Message, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	root, err := chf.NewRootPackedCHF(seg)
	if err != nil {
		return nil, err
	}

	msgs, err := root.NewMsgs(int32(len(flows)))
	if err != nil {
		return nil, err
	}

	for i, f := range flows {
		var list chf.Custom_List
		if n := int32(len(f.Customs)); n > 0 {
			if list, err = chf.NewCustom_List(seg, n); err != nil {
				return nil, err
			}
		}
		f.FillCHF(msgs.At(i), list)
	}

	root.SetMsgs(msgs)

	return msg, nil
}
