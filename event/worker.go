package event

import (
	"context"
	"log"
)

type worker struct {
	ep   *EventPool
	done chan bool
}

func newWorker(ep *EventPool) *worker {
	return &worker{
		ep:   ep,
		done: make(chan bool),
	}
}

func (w *worker) run() {
	log.Println("worker start")
	for {
		select {
		case <-w.done:
			log.Println("worker done")
			return
		case e := <-w.ep.eventChan:
			if e == nil {
				continue
			}

			e.Rsp = e.Cmd.ExecuteWithContext(context.Background())
			close(e.done)
		}
	}
}

func (w *worker) stop() {
	w.done <- true
	close(w.done)
}
