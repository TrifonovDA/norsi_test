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

const query_delete_node = "delete from notes where id = $1;"

type delete_request struct {
	ID uint64 `json:"id"`
}
type delete_result struct {
	Status string `json:"status"`
}

func (srv *Server) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	errc := make(chan handle_errors.Errors, 1) // канал ошибок
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	go func() {
		var mu sync.Mutex
		mu.Lock()
		defer mu.Unlock()
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadRequest, Err: err}
			return
		}

		var req = delete_request{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadRequest, Err: err}
			return
		}

		srv.Cache.Delete(req.ID)
		_ = srv.Dbconn.QueryRow(context.Background(), query_delete_node, req.ID)
		if err != nil {
			errc <- handle_errors.Errors{Code: http.StatusBadGateway, Err: err}
			return
		}

		done <- true
	}()

	select {
	case <-ctx.Done():
		http.Error(w, fmt.Errorf("timeout").Error(), http.StatusGatewayTimeout)
	case err := <-errc: // обработка ошибок
		http.Error(w, err.Err.Error(), err.Code)
	case <-done:
		Resp := &delete_result{Status: "success"}
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
