package worker

import (
	"beanstock/internal/pq"
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	ticker *time.Ticker
	pq     pq.PriorityQueue
	sync.Mutex
	done chan bool
}

func NewScheduler() *Scheduler {
	scheduler := &Scheduler{
		pq:     make(pq.PriorityQueue, 0),
		done:   make(chan bool),
		ticker: time.NewTicker(1 * time.Minute),
	}

	heap.Init(&scheduler.pq)
	return scheduler
}

func (s *Scheduler) start() {
	go func() {
		for {
			select {
			case <-s.done:
				s.ticker.Stop()
				return
			case t := <-s.ticker.C:
				fmt.Println("Tick at ", t)
				s.Lock()
				if len(s.pq) != 0 {
					message := heap.Pop(&s.pq).(*pq.DelayedMessage)
					now := time.Now().Unix()
					if message.Priority-now > 0 {
						heap.Push(&s.pq, message)
					} else {
						go message.Value.Execute()
					}
				}
				s.Unlock()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.done <- true
}

func (s *Scheduler) QueueMessage(message pq.Message, delay time.Duration) {
	s.Lock()
	now := time.Now().Unix()
	timeUntilExecute := now + int64(delay.Seconds())
	delayedMessage := &pq.DelayedMessage{
		Value:    message,
		Priority: timeUntilExecute,
	}
	heap.Push(&s.pq, delayedMessage)
	s.Unlock()
}
