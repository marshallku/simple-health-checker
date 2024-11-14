package handler

import (
	"net/http"

	"github.com/marshallku/statusy/templates"
	"github.com/marshallku/statusy/types"
)

func (s *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	results := s.store.GetResults()
	resultsList := make([]types.CheckResult, 0, len(results))
	for _, result := range results {
		resultsList = append(resultsList, result)
	}
	templates.IndexTemplate.Execute(w, resultsList)
}
