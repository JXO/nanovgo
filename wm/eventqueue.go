package wm

import (
    "sync"
)

// EventQueue is an infinitely buffered double-ended queue of Events. The zero value
// is usable, but a Queue value must not be copied
type EventQueue struct {
    mu sync.Mutex
    cond sync.Cond // cond.L is lazily initialized to &EventQueue.mu
    back []interface{} // FIFO
    front []interface{} // LIFO
}

func (q *EventQueue) lockAndInit() {
    q.mu.Lock()
    if q.cond.L == nil {
        q.cond.L = &q.mu
    }
}

func (q *EventQueue) NextEvent() interface{} {
    q.lockAndInit()
    defer q.mu.Unlock()

    for {
        if n := len(q.front); n > 0 {
            e := q.front[n-1]
            q.front[n-1] = nil
            q.front = q.front[:n-1]
            return e
        }

        if n := len(q.back); n > 0 {
            e := q.back[0]
            q.back[0] = nil
            q.back = q.back[1:]
            return e
        }

        q.cond.Wait()
    }
}

func (q *EventQueue) Send(event interface{}) {
    q.lockAndInit()
    defer q.mu.Unlock()

    q.back = append(q.back, event)
    q.cond.Signal()
}

func (q *EventQueue) SendFirst(event interface{}) {
    q.lockAndInit()
    defer q.mu.Unlock()

    q.front = append(q.front, event)
    q.cond.Signal()
}
