package server

import "testing"
import "log"
import "os"
import "github.com/lukesmith/cimple/chore"
import "github.com/lukesmith/cimple/messages"
import "github.com/satori/go.uuid"

func Test_Perform_BuildGitRepositoryJob(t *testing.T) {
	fakeMessageSender := &fakeMessageSender{
		sent: make([]*messages.Envelope, 0),
	}
	agent := newAgent(uuid.NewV4(), fakeMessageSender, log.New(os.Stdout, "", 0))
	job := &buildGitRepositoryJob{
		Url:    "https://github.com/cimpleci/test",
		Commit: "master",
	}

	chore := &chore.Chore{
		Job: job,
	}

	go func() {
		agent.Perform(chore)
	}()

	agent.available <- true

	sentMessage := (fakeMessageSender.sent[0].Body).(*messages.BuildGitRepository)
	if sentMessage.Url != job.Url {
		t.Errorf("Expected envelope body to have Url, was %+v", sentMessage.Url)
	}

	if sentMessage.Commit != job.Commit {
		t.Errorf("Expected envelope body to have Url, was %+v", sentMessage.Commit)
	}
}

type fakeMessageSender struct {
	sent []*messages.Envelope
}

func (s *fakeMessageSender) Close() error {
	return nil
}

func (s *fakeMessageSender) ReadMessage(envelope *messages.Envelope) error {
	return nil
}

func (s *fakeMessageSender) SendMessage(envelope *messages.Envelope) error {
	s.sent = append(s.sent, envelope)
	return nil
}
