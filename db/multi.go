package db

import (
	"sync"
)

type inMemoryMTSafe struct {
	Storage

	// Lock the whole storage for read/write, but read write concurrently
	rw sync.RWMutex
}

func NewInMemoryMTSafeStorage() Storage {
	return &inMemoryMTSafe{
		Storage: NewInMemoryStorage(),
	}
}

func (s *inMemoryMTSafe) Load(name UserId) (*User, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	return s.Storage.Load(name)
}

func (s *inMemoryMTSafe) Save(user *User) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.Storage.Save(user)
}

func (s *inMemoryMTSafe) Transaction() Tx {
	return s
}

func (s *inMemoryMTSafe) PutIfAbsent(user *User) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.Storage.Transaction().PutIfAbsent(user)
}

func (s *inMemoryMTSafe) Change(name UserId, val int64, expected Purse) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	return s.Storage.Transaction().Change(name, val, expected)
}
