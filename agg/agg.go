package agg

import (
	"sync"
	"time"

	capnp "zombiezen.com/go/capnproto2"

	"github.com/kentik/libkflow/flow"
	"github.com/kentik/libkflow/metrics"
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
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		a.error(err)
		return
	}

	a.Lock()
	a.metrics.TotalFlowsIn.Mark(a.queued)
	a.queued = 0
	flows, count, resampleRateAdj := a.queue.Dequeue(a.batchSize, a.batchSize)
	a.Unlock()

	if count == 0 {
		return
	}

	// adjust the sample rate for the provided flows
	normalizeSampleRate(flows, resampleRateAdj)

	// serialize the data using the provided segment (backed by msg)
	message, err := flow.ToCapnProtoMessage(flows, seg)
	if err != nil {
		a.error(err)
		return
	}

	a.output <- message

	a.metrics.TotalFlowsOut.Mark(int64(count))
}

func (a *Agg) error(err error) {
	select {
	case a.errors <- err:
	default:
	}
}

// normalizeSampleRate adjusts the sample rate in place on the provided [flow.Flow] slice based on a provided
// adjustment factor if it is > 1.0. The adjustment factor is multiplied by the original sample rate and 100 to get
// the new sample rate, as it is expected that a [flow.Flow] with a sample rate that does not account for this change.
func normalizeSampleRate(flows []flow.Flow, resampleRateAdj float32) {
	for i := range flows {
		sampleRate := flows[i].SampleRate
		adjustedSR := sampleRate * 100

		if resampleRateAdj > 1.0 {
			adjustedSR = uint32(float32(adjustedSR) * resampleRateAdj)
		}

		flows[i].SampleAdj = true
		flows[i].SampleRate = adjustedSR
	}
}
