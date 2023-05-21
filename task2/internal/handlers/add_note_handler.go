package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"simple_test/task2/pkg/handle_errors"
	"sync"
	"time"
)

type insert_result struct {
	Status string `json:"status"`
}

const insert_new_note = "insert into public.notes(id, title) values ($1, $2);"

func (srv *Server) AddNoteHandler(w http.ResponseWriter, req *http.Request) {
	done := make(chan bool)
	errc := make(chan handle_errors.Errors, 1) // канал ошибок
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	go func() {
		var mu sync.Mutex
		mu.Lock()
		defer mu.Unlock()

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadRequest, Err: err}
			return
		}

		var note = Note{}
		err = json.Unmarshal(body, &note)
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadRequest, Err: err}
			return
		}

		row, err := srv.Dbconn.Query(context.Background(), insert_new_note, note.ID, note.Title)
		defer row.Close()
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadGateway, Err: err}
			return
		}

		log.Println(note)
		srv.Cache.Insert(note)

		done <- true
	}()

	select {
	case <-ctx.Done():
		http.Error(w, fmt.Errorf("timeout").Error(), http.StatusGatewayTimeout)
	case err := <-errc: // обработка ошибок
		http.Error(w, err.Err.Error(), err.Code)
	case <-done:
		Resp := &insert_result{Status: "success"}
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
			http.Error(w, fmt.Errorf("sending response error").Error(), http.StatusInternalServerError)
			return
		}
	}
	return
}
