package chore

import (
	"log"
	"testing"
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

	chore := &Chore{}
	err := pool.QueueChore(chore)
	if err != nil {
		t.Fatalf("Expected no error queuing chore %s", err)
	}

	if chore.Completed {
		t.Fatalf("Expected chore not to have been completed")
	}
}

func TestWorkPool_QueueChore_WithWorker(t *testing.T) {
	pool := NewWorkPool()
	pool.AddWorker(&SampleWorker{})

	chore := &Chore{}
	err := pool.QueueChore(chore)
	if err != nil {
		t.Fatalf("Expected no error queuing chore %s", err)
	}

	if !chore.Completed {
		t.Fatalf("Expected chore to have been completed")
	}
}

type SampleWorker struct {
	ID   int32
	busy bool
}

func (worker *SampleWorker) CanPerform(c *Chore) bool {
	log.Printf("Checking if can perform chore %d on %s", c.ID, worker.ID)
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
