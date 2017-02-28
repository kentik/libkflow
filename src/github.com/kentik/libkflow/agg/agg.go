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
	interval  time.Duration
	ticker    *time.Ticker
	queue     *Queue
	queued    int64
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

	a := &Agg{
		output:   make(chan *capnp.Message),
		done:     make(chan struct{}),
		errors:   make(chan error, 100),
		interval: interval,
		ticker:   time.NewTicker(interval),
		msg:      msg,
		seg:      seg,
		metrics:  metrics,
	}

	a.Configure(fps)
	go a.aggregate()

	return a, nil
}

func (a *Agg) Configure(fps int) {
	var (
		interval_ms = float32(a.interval / time.Millisecond)
		batchSize   = (float32(fps) / 1000.0) * interval_ms
		buffer      = (float32(MaxFlowBuffer*fps) / 1000.0) * interval_ms
	)

	a.Lock()
	a.queue = New(int(buffer))
	a.batchSize = int(batchSize)
	a.Unlock()
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
	a.Lock()
	a.queued++
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
			close(a.output)
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
	a.metrics.TotalFlowsIn.Mark(a.queued)
	a.queued = 0
	flows, count, resampleRateAdj := a.queue.Dequeue(a.batchSize, a.batchSize)
	a.msg, msg = msg, a.msg
	a.seg, seg = seg, a.seg
	a.Unlock()

	if count == 0 {
		return
	}

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

	for i, f := range flows {
		sampleRate = f.SampleRate()
		adjustedSR = sampleRate * 100

		if resampleRateAdj > 1.0 {
			adjustedSR = uint32(float32(adjustedSR) * resampleRateAdj)
		}

		f.SetSampleAdj(true)
		f.SetSampleRate(adjustedSR)

		msgs.Set(i, *f)
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
