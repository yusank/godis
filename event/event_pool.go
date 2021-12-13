package event

type EventPool struct {
	eventChan chan *Event
	workers   []*worker
}

func NewEventPool() *EventPool {
	ep := &EventPool{
		eventChan: make(chan *Event, 1),
	}

	ep.startWorker(1)
	return ep
}

func (ep *EventPool) AddEvent(e *Event) {
	ep.eventChan <- e
}

func (ep *EventPool) Stop() {
	for _, w := range ep.workers {
		w.stop()
	}
	close(ep.eventChan)
}

func (ep *EventPool) startWorker(n int) {
	if n < 1 {
		n = 1
	}

	ep.workers = make([]*worker, n)
	for i := 0; i < n; i++ {
		ep.workers[i] = newWorker(ep)
		go ep.workers[i].run()
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
