package breaker

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	metric2 "github.com/1024casts/snake/pkg/metric"

	"github.com/1024casts/snake/pkg/errcode"
	"github.com/1024casts/snake/pkg/log"
)

// sreBreaker is a sre CircuitBreaker pattern.
type sreBreaker struct {
	stat metric2.RollingCounter
	r    *rand.Rand
	// rand.New(...) returns a non thread safe object
	randLock sync.Mutex

	k       float64
	request int64

	state int32
}

func newSRE(c *Config) Breaker {
	counterOpts := metric2.RollingCounterOpts{
		Size:           c.Bucket,
		BucketDuration: time.Duration(int64(c.Window) / int64(c.Bucket)),
	}
	stat := metric2.NewRollingCounter(counterOpts)
	return &sreBreaker{
		stat: stat,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),

		request: c.Request,
		k:       c.K,
		state:   StateClosed,
	}
}

func (b *sreBreaker) summary() (success int64, total int64) {
	b.stat.Reduce(func(iterator metric2.Iterator) float64 {
		for iterator.Next() {
			bucket := iterator.Bucket()
			total += bucket.Count
			for _, p := range bucket.Points {
				success += int64(p)
			}
		}
		return 0
	})
	return
}

func (b *sreBreaker) Allow() error {
	success, total := b.summary()
	k := b.k * float64(success)

	log.Info("breaker: request: %d, succee: %d, fail: %d", total, success, total-success)

	// check overflow requests = K * success
	if total < b.request || float64(total) < k {
		if atomic.LoadInt32(&b.state) == StateOpen {
			atomic.CompareAndSwapInt32(&b.state, StateOpen, StateClosed)
		}
		return nil
	}
	if atomic.LoadInt32(&b.state) == StateClosed {
		atomic.CompareAndSwapInt32(&b.state, StateClosed, StateOpen)
	}
	dr := math.Max(0, (float64(total)-k)/float64(total+1))
	drop := b.trueOnProba(dr)

	log.Info("breaker: drop ratio: %f, drop: %t", dr, drop)

	if drop {
		return errcode.ErrServiceUnavailable
	}
	return nil
}

func (b *sreBreaker) MarkSuccess() {
	b.stat.Add(1)
}

func (b *sreBreaker) MarkFailed() {
	// NOTE: when client reject requets locally, continue add counter let the
	// drop ratio higher.
	b.stat.Add(0)
}

func (b *sreBreaker) trueOnProba(proba float64) (truth bool) {
	b.randLock.Lock()
	truth = b.r.Float64() < proba
	b.randLock.Unlock()
	return
}
