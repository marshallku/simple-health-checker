package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/marshallku/statusy/store"
	"github.com/marshallku/statusy/templates"
	"github.com/marshallku/statusy/types"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Be more restrictive in production
	},
}

type Server struct {
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received")
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

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.store.AddClient(conn)
	defer s.store.RemoveClient(conn)

	// Send initial data
	conn.WriteJSON(store.Message{
		Type: "results",
		Data: s.store.GetResults(),
	})

	// Keep connection alive
	for {
		// Read message (required by WebSocket protocol)
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.HandleIndex)
	http.HandleFunc("/history", s.HandleHistory)
	http.HandleFunc("/ws", s.HandleWebSocket)
	fmt.Println("Server started on http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}
