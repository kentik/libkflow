package agg

import (
	"sync"
	"time"

	"github.com/kentik/libkflow/chf"
	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/metrics"
	capnp "zombiezen.com/go/capnproto2"
)

type Agg struct {
	output    chan *capnp.Message
	quitting  chan struct{} // closed when we're shutting down
	done      chan struct{} // closed when we're all done
	errors    chan error
	interval  time.Duration
	ticker    *time.Ticker
	queue     *Queue
	queued    int64
	batchSize int
	metrics   *metrics.Metrics
	sync.RWMutex
}

// MaxFlowBuffer defines the maximum amount of time in seconds to
// buffer flows at the maximum rate.
const MaxFlowBuffer = 8

// NewAgg creates a new Agg that aggregates flows into a single
// cap'n proto message after the specified interval, resampling
// as necessary to keep the total number under the fps arg.
func NewAgg(interval time.Duration, fps int, metrics *metrics.Metrics) (*Agg, error) {
	a := &Agg{
		output:   make(chan *capnp.Message),
		quitting: make(chan struct{}),
		done:     make(chan struct{}),
		errors:   make(chan error, 100),
		interval: interval,
		ticker:   time.NewTicker(interval),
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
	select {
	case <-a.quitting:
		// already closed
	default:
		close(a.quitting)
	}
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

func (a *Agg) Add(flow *flow.Flow) {
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
		case <-a.quitting:
			a.dispatch()
			close(a.output)
			close(a.done)
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

	a.metrics.BufferSize.Update(int64(a.queue.Length()))

	a.Lock()
	a.metrics.TotalFlowsIn.Mark(a.queued)
	a.queued = 0
	flows, count, resampleRateAdj := a.queue.Dequeue(a.batchSize, a.batchSize)
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
		sampleRate = f.SampleRate
		adjustedSR = sampleRate * 100

		if resampleRateAdj > 1.0 {
			adjustedSR = uint32(float32(adjustedSR) * resampleRateAdj)
		}

		f.SampleAdj = true
		f.SampleRate = adjustedSR

		var list chf.Custom_List
		if n := int32(len(f.Customs)); n > 0 {
			if list, err = chf.NewCustom_List(seg, n); err != nil {
				a.error(err)
				return
			}
		}

		f.FillCHF(msgs.At(i), list)
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
