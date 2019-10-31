package rotating

import (
	"sync"
	"time"
)

type Counter struct {
	blocks         []uint64
	total          uint64
	size           uint
	head           uint
	quitSignal     chan bool
	buffer         chan uint64
	blocking       bool
	blockingWait   sync.WaitGroup
	rotationTicker *time.Ticker
}

// NewCounter returns a new Rotating Counter that will keeps data only in a short period
// This period can be specified by the duration arguments.
// Counter works approximately and you can change accuracy by resolution argument.
// If resolution is too small Counter will use more memory and will be slower.
// It reserves duration/resolution blocks. So if you specified duration to 60 seconds and resolution to 500 milliseconds
// 120 blocks will be be reserved.
// With specifying bufferSize to zero Add function will wait until change is done but If you set bufferSize to 1 counter
// is non blocking and eventually consistent which would be faster. Setting bufferSize to bigger number will suitable
// for concurrency.
// Stop the counter to release associated resources.
func NewCounter(duration, resolution time.Duration, bufferSize uint8) *Counter {
	size := uint(duration / resolution)
	resolution = time.Duration(uint(duration) / size)

	counter := &Counter{
		blocks:         make([]uint64, size),
		total:          0,
		size:           size,
		quitSignal:     make(chan bool),
		buffer:         make(chan uint64, bufferSize),
		blocking:       bufferSize == 0,
		rotationTicker: time.NewTicker(resolution),
	}
	counter.start()
	return counter
}

func (c *Counter) Total() uint64 {
	return c.total
}

func (c *Counter) Add(value uint64) {
	if c.blocking {
		c.blockingWait.Add(1)
		defer c.blockingWait.Wait()
	}

	c.buffer <- value
}

func (c *Counter) start() {
	go func() {
		for {
			select {
			case value := <-c.buffer:
				c.total += value
				c.blocks[c.head] += value
				if c.blocking {
					c.blockingWait.Done()
				}

			case <-c.rotationTicker.C:
				c.head = (c.head + 1) % c.size
				c.total -= c.blocks[c.head]

				c.blocks[c.head] = 0

			case <-c.quitSignal:
				c.rotationTicker.Stop()
				close(c.quitSignal)
				close(c.buffer)
				return
			}
		}
	}()
}

func (c *Counter) Stop() {
	c.quitSignal <- true
}

func (c *Counter) Flush()  {
	c.blocks = make([]uint64, len(c.blocks))
	c.total = 0
}
