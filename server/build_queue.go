package server

import (
	"github.com/lukesmith/cimple/chore"
	"log"
)

type BuildQueue interface {
	Queue(job interface{})
}

type buildQueue struct {
	queue     chan interface{}
	agentpool *agentpool
}

func (bq *buildQueue) Queue(job interface{}) {
	bq.queue <- job
}

func (a *buildQueue) run() {
	log.Print("Running....")
	for {
		select {
		case i := <-a.queue:
			chore := &chore.Chore{
				Done: make(chan bool),
				Job:  i,
			}
			a.agentpool.workerpool.QueueChore(chore)
			go func() {
				<-chore.Done
			}()
		}
	}
}
