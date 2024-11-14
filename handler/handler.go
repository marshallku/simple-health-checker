package handler

import (
	"fmt"
	"net/http"

	"github.com/marshallku/statusy/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (s *Handler) RegisterRoutes() error {
	http.HandleFunc("/", s.HandleIndex)
	http.HandleFunc("/history", s.HandleHistory)
	http.HandleFunc("/ws", s.HandleWebSocket)

	fmt.Println("Server started on http://localhost:8080")

	return http.ListenAndServe(":8080", nil)
}
