package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/messages"
	"log"
	"net/http"
)

type hooksHandler struct {
	agentpool *agentpool
	router    *mux.Router
}

func (h *hooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHooks(agents *agentpool) http.Handler {
	router := mux.NewRouter()

	handler := &hooksHandler{
		router:    router,
		agentpool: agents,
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

		for k, _ := range h.agentpool.agents {
			k.send(&msg)
		}

		w.Write([]byte("Accepted"))
	}
}
