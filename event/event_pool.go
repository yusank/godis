package event

import (
	"sync"
)

type EventPool struct {
	eventChan chan *Event
	workers   []*worker
	wg        *sync.WaitGroup
}

func NewEventPool() *EventPool {
	ep := &EventPool{
		eventChan: make(chan *Event, 1),
		wg:        new(sync.WaitGroup),
	}

	ep.startWorker(1)
	return ep
}

func (ep *EventPool) AddEvent(e *Event) {
	ep.eventChan <- e
}

func (ep *EventPool) Stop() {
	close(ep.eventChan)
	for _, w := range ep.workers {
		w.stop()
	}

	ep.wg.Wait()
}

func (ep *EventPool) startWorker(n int) {
	if n < 1 {
		n = 1
	}

	ep.workers = make([]*worker, n)
	for i := 0; i < n; i++ {
		ep.workers[i] = newWorker(ep)
		ep.wg.Add(1)
		go ep.workers[i].run(ep.wg)
	}
}

var globalEventPool *EventPool

func init() {
	SetGlobalEventPool(NewEventPool())
}

func SetGlobalEventPool(ep *EventPool) {
	globalEventPool = ep
}

func AddEvent(e *Event) {
	globalEventPool.AddEvent(e)
}

func Stop() {
	globalEventPool.Stop()
}
