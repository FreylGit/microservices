package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	baseUrl       = "localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"isPublic"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type SyncMap struct {
	elems map[int64]*Note
	m     sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to decode note data", http.StatusBadRequest)
		return
	}
	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	note := &Note{
		ID:        rand.Int63(),
		Info:      *info,
		CreatedAt: now,
		UpdateAt:  now,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
	notes.m.Lock()
	defer notes.m.Unlock()
	notes.elems[note.ID] = note
}
func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "id")
	id, err := parseNoteId(noteId)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusNotFound)
		return
	}
	notes.m.RLock()
	defer notes.m.RUnlock()
	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "Invalid note ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Invalid note ID", http.StatusInternalServerError)
		return
	}

}
func main() {
	r := chi.NewRouter()
	r.Post(createPostfix, createNoteHandler)
	r.Get(getPostfix, getNoteHandler)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}

func parseNoteId(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}