package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/chore"
	"github.com/lukesmith/cimple/messages"
	"log"
	"net/http"
)

type hooksHandler struct {
	router     *mux.Router
	buildQueue *buildQueue
}

func (h *hooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHooks(buildQueue *buildQueue) http.Handler {
	router := mux.NewRouter()

	handler := &hooksHandler{
		router:     router,
		buildQueue: buildQueue,
	}

	router.HandleFunc("/hooks", handler.postHook).
		Methods("POST")

	return handler
}

func (h *hooksHandler) postHook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var msg messages.BuildGitRepository
	err := decoder.Decode(&msg)
	if err != nil {
		log.Printf("Failed to understand build hook - %+v", err)

		w.Write([]byte("Unprocessible entity"))
	} else {

		log.Printf("Received build hook for %s, %s", msg.Url, msg.Commit)

		h.buildQueue.queue <- &msg

		w.Write([]byte("Accepted"))
	}
}

type buildQueue struct {
	queue     chan interface{}
	agentpool *agentpool
}

func (a *buildQueue) run() {
	log.Print("Running....")
	for {
		select {
		case i := <-a.queue:
			chore := &chore.Chore{}
			a.agentpool.workerpool.QueueChore(chore)
			log.Printf("Queued %s", i)
			for k, _ := range a.agentpool.agents {
				k.send(i)
			}
		}
	}
}
