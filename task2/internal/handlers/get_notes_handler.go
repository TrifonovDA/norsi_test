package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type get_notes_result struct {
	Status string `json:"status"`
	Notes  []Note `json:"notes"`
}

func (srv *Server) GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	Resp := &get_notes_result{}
	done := make(chan bool)
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)

	defer cancel()

	go func() {
		var mu sync.Mutex
		mu.Lock()
		defer mu.Unlock()
		notes := srv.Cache.Get_notes()
		Resp = &get_notes_result{Status: "success", Notes: notes}
		done <- true
	}()

	select {
	case <-ctx.Done():
		http.Error(w, fmt.Errorf("timeout").Error(), http.StatusGatewayTimeout)
	case <-done:
		Response_json, err := json.Marshal(Resp)
		if err != nil {
			log.Println("marshlling error:", err)
			http.Error(w, fmt.Errorf("marshalling response error").Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(Response_json)
		if err != nil {
			log.Println("sending result error:", err)
			http.Error(w, fmt.Errorf("marshalling response error").Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	return
}
