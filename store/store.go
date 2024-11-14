package store

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/marshallku/statusy/types"
)

type Store struct {
	mu        sync.RWMutex
	results   map[string]types.CheckResult
	history   []types.History
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewStore() *Store {
	s := &Store{
		results:   make(map[string]types.CheckResult),
		history:   make([]types.History, 0),
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
	}
	go s.handleBroadcast()
	return s
}

func (s *Store) handleBroadcast() {
	for message := range s.broadcast {
		s.mu.RLock()
		for client := range s.clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mu.RUnlock()
	}
}

func (s *Store) AddClient(client *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = true
}

func (s *Store) RemoveClient(client *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client)
}

func (s *Store) UpdateResult(result types.CheckResult) {
	s.mu.Lock()
	s.results[result.URL] = result
	s.mu.Unlock()

	s.broadcast <- Message{Type: "results", Data: s.GetResults()}

	// Add history
	status := "UP"
	if !result.Status {
		status = "DOWN"
	}

	s.AddHistory(types.History{
		URL:    result.URL,
		Status: status,
	})
}

func (s *Store) AddHistory(h types.History) {
	s.mu.Lock()
	s.history = append([]types.History{h}, s.history...)
	if len(s.history) > 10 {
		s.history = s.history[:10]
	}
	s.mu.Unlock()

	s.broadcast <- Message{Type: "history", Data: s.GetHistory()}
}

func (s *Store) GetResults() map[string]types.CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.results
}

func (s *Store) GetHistory() []types.History {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.history
}
