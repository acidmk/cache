package lru

import (
	"sync/atomic"
	"time"
)

type clock interface {
	Now() time.Time
	Stop()
}

// ClockNone fake clock, used for cache without expiration
type ClockNone struct{}

func (c *ClockNone) Now() time.Time {
	return time.Time{}
}
func (c *ClockNone) Stop() {}

func newClockNone() clock {
	return &ClockNone{}
}

// ClockSimple is a default clock. Used for precise expiration
type ClockSimple struct{}

func (c *ClockSimple) Now() time.Time {
	return time.Now()
}
func (c *ClockSimple) Stop() {}

func newClockSimple() clock {
	return &ClockSimple{}
}

// ClockDiscrete is an optimized clock. Not as precise as ClockSimple but significantly faster.
// Current time is refreshed each 500ms.
type ClockDiscrete struct {
	updateTicker *time.Ticker
	value        *atomic.Value
}

func (c *ClockDiscrete) Now() time.Time {
	now := c.value.Load().(time.Time)
	return now
}

func (c *ClockDiscrete) Stop() {
	c.updateTicker.Stop()
}

func (c *ClockDiscrete) refresh() {
	c.value.Store(time.Now())
}

func newClockDiscrete(updateTime time.Duration) clock {
	if updateTime == 0 {
		updateTime = time.Millisecond * 500
	}

	ret := &ClockDiscrete{
		updateTicker: time.NewTicker(updateTime),
		value:        &atomic.Value{},
	}

	ret.refresh()

	go func(c *ClockDiscrete) {
		for range c.updateTicker.C {
			c.refresh()
		}
	}(ret)

	return ret
}
