package ringbuffer

import (
	"log"
	"sync"
	"time"

	"github.com/docker/docker/daemon/logger"
)

type RingBuffer struct {
	bfs  []logger.Message
	rbfs []logger.Message
	tk   *time.Ticker
	rl   sync.RWMutex

	head int
	iput int
	size int
}

func (rb *RingBuffer) Write(p logger.Message) (n int, err error) {
	rb.bfs[rb.iput] = p
	rb.iput += 1
	rb.iput %= rb.size

	return n, err
}

func (rb *RingBuffer) Tail(n int, Since time.Time) []*logger.Message {
	rb.rl.RLock()
	defer rb.rl.RUnlock()
	if rb.size < n || n <= 0 {
		n = rb.size
	}
	log.Println(`???????????`, n, rb.size)
	tmp := make([]*logger.Message, n)
	for i := 1; i <= n; i++ {
		id := (rb.head - i + rb.size) % rb.size
		insert := n - i
		if rb.rbfs[id].Timestamp.Before(Since) || len(rb.rbfs[id].Line) == 0 {
			tmp = tmp[insert+1:]

			break
		}

		tmp[insert] = &rb.rbfs[id]
		tmp[insert].Line = []byte(string(tmp[insert].Line))
		if id == rb.head {
			break
		}

	}
	return tmp
}

func (rb *RingBuffer) rwSplitting() {
	go func() {
		for range rb.tk.C {
			rb.rl.Lock()
			copy(rb.rbfs, rb.bfs)
			rb.head = rb.iput
			rb.rl.Unlock()
		}
	}()
}
func (rb *RingBuffer) Close() {
	rb.tk.Stop()
}
func New(size int) (rb *RingBuffer) {
	rb = &RingBuffer{
		bfs:  make([]logger.Message, size),
		rbfs: make([]logger.Message, size),
		size: size,
		rl:   sync.RWMutex{},
		tk:   time.NewTicker(time.Second),
	}
	rb.rwSplitting()
	return rb
}
