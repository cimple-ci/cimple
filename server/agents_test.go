package server

import (
	"fmt"
	"github.com/lukesmith/cimple/web_application"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	server     *httptest.Server
	agentsUrl  string
	reader     io.Reader
	app        *web_application.Application
	agentsPool AgentPool
)

type fakeAgentPool struct {
	agents []*Agent
}

func (a *fakeAgentPool) GetAgents() ([]*Agent, error) {
	return a.agents, nil
}

func (a *fakeAgentPool) Join(agent *Agent) {
	a.agents = append(a.agents, agent)
}

func (a *fakeAgentPool) Leave(agent *Agent) {
}

func init() {
	options := &web_application.ApplicationOptions{}
	app = web_application.NewApplication(options)
	agentsPool = &fakeAgentPool{}
	registerAgents(app, agentsPool, log.New(os.Stdout, "", 0))
	server = httptest.NewServer(app)
	agentsUrl = fmt.Sprintf("%s/agents", server.URL)
}

func Test_GetAgents_Json(t *testing.T) {
	agent := newAgent(uuid.NewV4(), nil, nil)
	agent.ConnectedDate = time.Date(2016, time.July, 20, 20, 12, 55, 826456124, time.Local)
	agent.busy = true
	agentsPool.Join(agent)

	request, err := http.NewRequest("GET", agentsUrl, reader)
	request.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
	}

	body, _ := ioutil.ReadAll(res.Body)
	actual := string(body)
	expected := fmt.Sprintf(`[{"Id":"%s","ConnectedDate":"2016-07-20T20:12:55.826456124+01:00","Busy":true}]`, agent.Id)

	if actual != expected {
		t.Errorf("Response not what was expected %s - %s", expected, actual)
	}
}
