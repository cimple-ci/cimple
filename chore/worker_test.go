package chore

import (
	"log"
	"testing"
	"time"
)

func TestWorkPool_AddWorker(t *testing.T) {
	pool := NewWorkPool()

	w := &SampleWorker{}
	err := pool.AddWorker(w)
	if err != nil {
		t.Fatalf("Expected no error adding worker - %s", err)
	}
}

func TestWorkPool_RemoveWorker(t *testing.T) {
	pool := NewWorkPool()
	w := &SampleWorker{}

	err := pool.RemoveWorker(w)
	if err != nil {
		t.Fatalf("Expected no error removing worker - %s", err)
	}
}

func TestWorkPool_QueueChore_WithoutWorkers(t *testing.T) {
	pool := NewWorkPool()

	done := make(chan bool)
	chore := &Chore{Done: done}
	err := pool.QueueChore(chore)
	if err != nil {
		t.Fatalf("Expected no error queuing chore %s", err)
	}

	select {
	case <-done:
		t.Fatalf("Expected chore not to have been completed")
	case <-time.After(5 * time.Second):
		t.Log("Chore not completed, as expected")
	}
}

func TestWorkPool_QueueChore_WithWorker(t *testing.T) {
	pool := NewWorkPool()
	pool.AddWorker(&SampleWorker{})

	done := make(chan bool)
	chore := &Chore{Done: done}
	err := pool.QueueChore(chore)
	if err != nil {
		t.Fatalf("Expected no error queuing chore %s", err)
	}

	select {
	case <-done:
		t.Log("Chore done")
	case <-time.After(10 * time.Second):
		t.Fatalf("Expected chore to have been completed")
	}
}

type SampleWorker struct {
	ID   int32
	busy bool
}

func (worker *SampleWorker) CanPerform(c *Chore) bool {
	log.Printf("Checking if can perform chore %d on %+v", c.ID, worker.ID)
	return !worker.busy
}

func (worker *SampleWorker) Perform(c *Chore) error {
	worker.busy = true
	defer func() {
		worker.busy = false
	}()
	log.Printf("Performing chore %d", c.ID)
	log.Printf("Performed chore %d", c.ID)
	return nil
}
