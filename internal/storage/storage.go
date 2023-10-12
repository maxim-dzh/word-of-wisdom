package storage

import (
	"sync"

	"github.com/maxim-dzh/word-of-wisdom/internal/hashcash"
)

type storage struct {
	sync.RWMutex
	challenges map[string]*hashcash.Header
}

// SetChallenge saves a challenge
func (s *storage) SetChallenge(id string, header *hashcash.Header) {
	s.Lock()
	s.challenges[id] = header
	s.Unlock()
}

// GetChallenge returns a challenge by id
func (s *storage) GetChallenge(id string) (header *hashcash.Header, exists bool) {
	s.RLock()
	header, exists = s.challenges[id]
	s.RUnlock()
	return
}

// DeleteChallenge deletes a challenge by id
func (s *storage) DeleteChallenge(id string) {
	s.Lock()
	delete(s.challenges, id)
	s.Unlock()
}

// NewStorage returns a new storage instance
func NewStorage() *storage {
	return &storage{
		challenges: make(map[string]*hashcash.Header),
	}
}
