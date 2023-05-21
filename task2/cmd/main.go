package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"simple_test/task2/internal/handlers"
	"simple_test/task2/pkg/db_tools"
	"simple_test/task2/pkg/logger"
)

func main() {
	dbconn := db_tools.NewConnection(context.Background())
	defer dbconn.Close()
	cache := handlers.Simple_cache{}
	cache = cache.NewCache()
	err := db_tools.Up_cache(context.Background(), &cache, dbconn)
	if err != nil {
		log.Println(err)
	}

	srv := handlers.Server{
		Dbconn: dbconn,
		Cache:  &cache,
	}

	r := mux.NewRouter()
	r.Use(logger.LogRequest)
	r.HandleFunc("/notes", srv.GetNotesHandler).Methods("GET")
	r.HandleFunc("/notes", srv.AddNoteHandler).Methods("POST")
	r.HandleFunc("/notes", srv.DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
