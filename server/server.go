package server

import (
	"net/http"

	"github.com/marshallku/statusy/store"
	"github.com/marshallku/statusy/templates"
	"github.com/marshallku/statusy/types"
)

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	results := s.store.GetResults()
	resultsList := make([]types.CheckResult, 0, len(results))
	for _, result := range results {
		resultsList = append(resultsList, result)
	}
	templates.IndexTemplate.Execute(w, resultsList)
}

func (s *Server) HandleHistory(w http.ResponseWriter, r *http.Request) {
	history := s.store.GetHistory()
	templates.HistoryTemplate.Execute(w, history)
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.HandleIndex)
	http.HandleFunc("/history", s.HandleHistory)
	return http.ListenAndServe(":8080", nil)
}
