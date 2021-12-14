package event

import (
	"log"
	"sync"
)

type worker struct {
	ep   *EventPool
	done chan bool
}

func newWorker(ep *EventPool) *worker {
	return &worker{
		ep:   ep,
		done: make(chan bool, 1),
	}
}

func (w *worker) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-w.done:
			log.Println("worker done")
			return
		case e := <-w.ep.eventChan:
			if e == nil {
				return
			}

			e.Rsp = <-e.Cmd.ExecuteAsync()
			close(e.done)
		}
	}
}

func (w *worker) stop() {
	w.done <- true
	close(w.done)
}
