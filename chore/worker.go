package chore

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
)

type Worker interface {
	CanPerform(c *Chore) bool
	Perform(c *Chore) error
}

type Chore struct {
	ID        int32
	Completed bool
}

type WorkPool struct {
	availableWorkers chan Worker
	removeWorker     chan Worker
	chores           chan *Chore
	active           int64
	workers          []Worker
	queuedChores     []*Chore
	checkTime        time.Duration
	busyWorkers      int64
}

func NewWorkPool() *WorkPool {
	workPool := &WorkPool{
		availableWorkers: make(chan Worker),
		removeWorker:     make(chan Worker),
		chores:           make(chan *Chore),
		workers:          make([]Worker, 0),
		queuedChores:     make([]*Chore, 0),
		checkTime:        1000 * time.Millisecond,
		busyWorkers:      0,
	}

	workPool.start()

	return workPool
}

func (wp *WorkPool) AddWorker(w Worker) error {
	wp.availableWorkers <- w
	return nil
}

func (wp *WorkPool) RemoveWorker(w Worker) error {
	wp.removeWorker <- w
	return nil
}

func (wp *WorkPool) QueueChore(c *Chore) error {
	wp.chores <- c
	return nil
}

func (wp *WorkPool) selectWorker(c *Chore) (Worker, error) {
	for _, w := range wp.workers {
		if w.CanPerform(c) {
			return w, nil
		}
	}

	return nil, errors.New("Unable to find a worker")
}

func (wp *WorkPool) run(c *Chore) (bool, error) {
	performed := false
	worker, err := wp.selectWorker(c)

	if worker != nil {
		err = worker.Perform(c)
		performed = true
		c.Completed = true
	} else {
		log.Printf("Unable to perform chore %d", c.ID)
	}

	return performed, err
}

func (wp *WorkPool) doChore(c *Chore) {
	atomic.AddInt64(&wp.busyWorkers, 1)
	ok, _ := wp.run(c)

	if !ok {
		wp.chores <- c
	}
	atomic.AddInt64(&wp.busyWorkers, -1)
}

func (wp *WorkPool) start() {
	go func() {
		timer := time.NewTimer(wp.checkTime)
		statTimer := time.NewTimer(wp.checkTime)

		for {
			select {
			case w := <-wp.availableWorkers:
				wp.workers = append(wp.workers, w)
			case w := <-wp.removeWorker:
				log.Printf("Removing worker %s", w)
			case c := <-wp.chores:
				wp.queuedChores = append(wp.queuedChores, c)
			case <-statTimer.C:
				chores := len(wp.queuedChores)
				workers := len(wp.workers)
				busyWorkers := atomic.LoadInt64(&wp.busyWorkers)
				log.Printf("WorkPool : Workers %d/%d : Chores %d", busyWorkers, workers, chores)
				statTimer.Reset(wp.checkTime)
			case <-timer.C:
				if len(wp.queuedChores) > 0 {
					c := &Chore{}
					c, wp.queuedChores = wp.queuedChores[len(wp.queuedChores)-1], wp.queuedChores[:len(wp.queuedChores)-1]
					go wp.doChore(c)
				}
				timer.Reset(wp.checkTime)
			}
		}
	}()
}
