package server

import (
	"github.com/lukesmith/cimple/chore"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type BuildQueue interface {
	Queue(job BuildJob)
	GetQueued() ([]BuildJob, error)
}

type BuildJob interface {
	Id() uuid.UUID
	SubmissionDate() time.Time
}

type buildGitRepositoryJob struct {
	id             uuid.UUID
	submissionDate time.Time
	Url            string
	Commit         string
}

func (bj *buildGitRepositoryJob) Id() uuid.UUID {
	return bj.id
}

func (bj *buildGitRepositoryJob) SubmissionDate() time.Time {
	return bj.submissionDate
}

func NewBuildGitRepositoryJob(url string, commit string) BuildJob {
	return &buildGitRepositoryJob{
		Url:            url,
		Commit:         commit,
		id:             uuid.NewV4(),
		submissionDate: time.Now(),
	}
}

type buildQueue struct {
	queue     chan interface{}
	agentpool *agentpool
}

func (bq *buildQueue) Queue(job BuildJob) {
	bq.queue <- job
}

func (bq *buildQueue) GetQueued() ([]BuildJob, error) {
	queued := make([]BuildJob, 0)

	chores, err := bq.agentpool.workerpool.QueuedChores()

	if err != nil {
		return nil, err
	}

	for _, chore := range chores {
		queued = append(queued, chore.Job.(BuildJob))
	}

	return queued, nil
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
