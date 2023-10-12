package wordofwisdom

import (
	"math/rand"
	"time"
)

type service struct {
	words []string
	rand  *rand.Rand
}

// GetWord returns a random word for the words list
func (s *service) GetWord() string {
	return s.words[s.rand.Intn(len(s.words))]
}

// NewService returns a new service instance
//
//nolint:gosec // it's for word randomizing
func NewService(words []string) *service {
	return &service{
		words: words,
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}
