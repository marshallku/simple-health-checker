package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/marshallku/statusy/store"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Be more restrictive in production
	},
}

func (s *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.store.AddClient(conn)
	defer s.store.RemoveClient(conn)

	conn.WriteJSON(store.Message{
		Type: "results",
		Data: s.store.GetResults(),
	})

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
