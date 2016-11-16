package agg

import (
	"sync"
	"time"

	"github.com/kentik/libkflow/chf"
	"zombiezen.com/go/capnproto2"
)

type Agg struct {
	output    chan *capnp.Message
	done      chan struct{}
	errors    chan error
	ticker    *time.Ticker
	queue     *Queue
	batchSize int
	msg       *capnp.Message
	seg       *capnp.Segment
	metrics   *Metrics
	sync.RWMutex
}

// MaxFlowBuffer defines the maximum amount of time in seconds to
// buffer flows at the maximum rate.
const MaxFlowBuffer = 8

// NewAgg creates a new Agg that aggregates flows into a single
// cap'n proto message after the specified interval, resampling
// as necessary to keep the total number under the fps arg.
func NewAgg(interval time.Duration, fps int, metrics *Metrics) (*Agg, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	batchSize := int((float32(fps) / 1000.0) * float32(interval))

	a := &Agg{
		output:    make(chan *capnp.Message),
		done:      make(chan struct{}),
		errors:    make(chan error, 100),
		ticker:    time.NewTicker(interval),
		queue:     New(MaxFlowBuffer * fps),
		batchSize: batchSize,
		msg:       msg,
		seg:       seg,
		metrics:   metrics,
	}

	go a.aggregate()

	return a, nil
}

func (a *Agg) Stop() {
	a.done <- struct{}{}
}

func (a *Agg) Output() <-chan *capnp.Message {
	return a.output
}

func (a *Agg) Done() <-chan struct{} {
	return a.done
}

func (a *Agg) Errors() <-chan error {
	return a.errors
}

func (a *Agg) Segment() *capnp.Segment {
	a.RLock()
	seg := a.seg
	a.RUnlock()
	return seg
}

func (a *Agg) Add(flow *chf.CHF) {
	a.metrics.TotalFlowsIn.Mark(1)
	a.Lock()
	if a.queue.Enqueue(flow) != nil {
		a.metrics.RateLimitDrops.Mark(1)
	}
	a.Unlock()
}

func (a *Agg) aggregate() {
	for {
		select {
		case <-a.ticker.C:
			a.dispatch()
		case <-a.done:
			a.dispatch()
			a.done <- struct{}{}
			return
		}
	}
}

func (a *Agg) dispatch() {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		a.error(err)
		return
	}

	a.Lock()
	flows, count, resampleRateAdj := a.queue.Dequeue(a.batchSize, a.batchSize)
	a.msg, msg = msg, a.msg
	a.seg, seg = seg, a.seg
	a.Unlock()

	root, err := chf.NewRootPackedCHF(seg)
	if err != nil {
		a.error(err)
		return
	}

	msgs, err := root.NewMsgs(int32(len(flows)))
	if err != nil {
		a.error(err)
		return
	}

	var sampleRate uint32
	var adjustedSR uint32
	count = 0

	for _, f := range flows {
		sampleRate = f.SampleRate()
		adjustedSR = sampleRate * 100

		if resampleRateAdj > 1.0 {
			adjustedSR = uint32(float32(adjustedSR) * resampleRateAdj)
		}

		f.SetSampleAdj(true)
		f.SetSampleRate(adjustedSR)
		msgs.Set(count, *f)

		count++
	}

	root.SetMsgs(msgs)
	a.output <- msg

	a.metrics.OrigSampleRate.Update(int64(sampleRate))
	a.metrics.NewSampleRate.Update(int64(adjustedSR))
	a.metrics.TotalFlowsOut.Mark(int64(count))
}

func (a *Agg) error(err error) {
	select {
	case a.errors <- err:
	default:
	}
}
