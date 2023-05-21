package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type Server struct {
	Dbconn *pgxpool.Pool
	Cache  *Simple_cache
	Ctx    context.Context
}

type Note struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type Simple_cache struct {
	mu_cache  sync.Mutex
	Cache_map map[uint64]string
}

func (cache *Simple_cache) NewCache() Simple_cache {
	return Simple_cache{
		mu_cache:  sync.Mutex{},
		Cache_map: make(map[uint64]string),
	}
}

func (cache *Simple_cache) Insert(note Note) {
	cache.mu_cache.Lock()
	defer cache.mu_cache.Unlock()

	cache.Cache_map[note.ID] = note.Title
}

func (cache *Simple_cache) Get_notes() (notes []Note) {
	cache.mu_cache.Lock()
	defer cache.mu_cache.Unlock()

	//notes = make([]Note, len(cache.Cache_map))

	for key, value := range cache.Cache_map {
		notes = append(notes, Note{
			ID:    key,
			Title: value,
		})
	}
	return notes
}

func (cache *Simple_cache) Delete(key uint64) {
	cache.mu_cache.Lock()
	defer cache.mu_cache.Unlock()

	delete(cache.Cache_map, key)
	return
}
