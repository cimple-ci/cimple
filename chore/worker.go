package chore

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Worker interface {
	CanPerform(c *Chore) bool
	Perform(c *Chore) error
}

type Chore struct {
	ID   int32
	Job  interface{}
	Done chan bool
}

type WorkPool struct {
	availableWorkers chan Worker
	removeWorker     chan Worker
	chores           chan *Chore
	active           int32
	workers          []Worker
	queuedChores     []*Chore
	checkTime        time.Duration
	busyWorkers      int32
	totalChores      int32
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
		totalChores:      0,
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
	atomic.AddInt32(&wp.totalChores, 1)
	c.ID = atomic.LoadInt32(&wp.totalChores)
	wp.chores <- c
	return nil
}

func (wp *WorkPool) QueuedChores() ([]*Chore, error) {
	return wp.queuedChores, nil
}

func (wp *WorkPool) selectWorker(c *Chore) (Worker, error) {
	for _, w := range wp.workers {
		if w.CanPerform(c) {
			return w, nil
		}
	}

	return nil, errors.New("Unable to find a worker")
}

func (wp *WorkPool) perform(c *Chore) error {
	performed := false
	worker, err := wp.selectWorker(c)
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		if worker != nil {
			atomic.AddInt32(&wp.busyWorkers, 1)
			err = worker.Perform(c)
			log.Printf("Chore %+v has been performed", c)
			performed = true
			atomic.AddInt32(&wp.busyWorkers, -1)
			c.Done <- performed
		} else {
			log.Printf("Unable to perform chore %+v", c)
			c.Done <- performed
		}
	}()

	wg.Wait()

	if !performed {
		log.Printf("Chore %+v did not successfully complete, requeuing", c)
		wp.chores <- c
	}

	return err
}

func (wp *WorkPool) start() {
	go func() {
		timer := time.NewTimer(wp.checkTime)
		statTimer := time.NewTimer(wp.checkTime)

		for {
			select {
			case w := <-wp.availableWorkers:
				log.Printf("Adding worker %s", w)
				wp.workers = append(wp.workers, w)
			case w := <-wp.removeWorker:
				log.Printf("Removing worker %s", w)
			case c := <-wp.chores:
				log.Printf("Chore %+v queued", c)
				wp.queuedChores = append(wp.queuedChores, c)
			case <-statTimer.C:
				chores := len(wp.queuedChores)
				workers := len(wp.workers)
				busyWorkers := atomic.LoadInt32(&wp.busyWorkers)
				log.Printf("WorkPool : Workers %d/%d : Chores %d", busyWorkers, workers, chores)
				statTimer.Reset(wp.checkTime)
			case <-timer.C:
				log.Printf("Checking queued chores for work")
				workers := len(wp.workers)
				busyWorkers := atomic.LoadInt32(&wp.busyWorkers)
				if len(wp.queuedChores) > 0 {
					if int(busyWorkers) < workers {
						c := &Chore{}
						c, wp.queuedChores = wp.queuedChores[len(wp.queuedChores)-1], wp.queuedChores[:len(wp.queuedChores)-1]
						log.Printf("Chore %d to be worked on", c.ID)
						go wp.perform(c)
					} else {
						log.Printf("No available workers. Workers %d/%d", busyWorkers, workers)
					}
				} else {
					log.Printf("No chores queued.")
				}
				timer.Reset(wp.checkTime)
			}
		}
	}()
}
