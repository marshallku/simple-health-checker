package handler

import (
	"net/http"

	"github.com/marshallku/statusy/templates"
)

func (s *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	history := s.store.GetHistory()
	templates.HistoryTemplate.Execute(w, history)
}
