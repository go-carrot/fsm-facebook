package fsmfacebook

import (
	"errors"

	"github.com/go-carrot/fsm"
)

type CacheStore struct {
	Traversers map[string]fsm.Traverser
}

func (s *CacheStore) FetchTraverser(uuid string) (fsm.Traverser, error) {
	if traverser, ok := s.Traversers[uuid]; ok {
		return traverser, nil
	}
	return nil, errors.New("Traverser does not exist")
}

func (s *CacheStore) CreateTraverser(uuid string) (fsm.Traverser, error) {
	if _, ok := s.Traversers[uuid]; ok {
		return nil, errors.New("Traverser with UUID already exists")
	}
	traverser := &CachedTraverser{}
	traverser.SetUUID(uuid)
	s.Traversers[uuid] = traverser
	return traverser, nil
}
