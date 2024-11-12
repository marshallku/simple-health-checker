package store

import (
	"sync"

	"github.com/marshallku/statusy/types"
)

type Store struct {
	mu      sync.RWMutex
	results map[string]types.CheckResult
	history []types.History
}

func NewStore() *Store {
	return &Store{
		results: make(map[string]types.CheckResult),
		history: make([]types.History, 0),
	}
}

func (s *Store) UpdateResult(result types.CheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results[result.URL] = result
}

func (s *Store) AddHistory(h types.History) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append([]types.History{h}, s.history...)
	if len(s.history) > 10 {
		s.history = s.history[:10]
	}
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
