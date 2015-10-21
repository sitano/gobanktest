package db

import "fmt"

type inMemory struct {
	data map[UserId]Purse
}

func NewInMemoryStorage() Storage {
	return &inMemory{
		data: map[UserId]Purse{},
	}
}

func (s *inMemory) Load(name UserId) (*User, error) {
	purse, ok := s.data[name]
	if !ok {
		return nil, fmt.Errorf("There is no such user %s", name)
	}

	return &User{
		Name: name,
		Purse: purse,
	}, nil
}

func (s *inMemory) Save(user *User) error {
	s.data[user.Name] = user.Purse
	return nil
}

func (s *inMemory) List() map[UserId]Purse {
	copy := map[UserId]Purse{}

	for name, purse := range s.data {
		copy[name] = purse
	}

	return copy
}

func (s *inMemory) Transaction() Tx {
	return s
}

// Single threaded implementation has atomicity guarantee by definition
func (s *inMemory) PutIfAbsent(user *User) error {
	if _, ok := s.data[user.Name]; ok {
		return fmt.Errorf("Can't put user %s into the storage: already present", user.Name)
	}

	return s.Save(user)
}

// Single threaded implementation has atomicity guarantee by definition
func (s *inMemory) Change(name UserId, val int64, expected Purse) error {
	user, err := s.Load(name)
	if err != nil {
		return err
	}

	if user.Purse != expected {
		return fmt.Errorf("User balance have been changed since. user.Purse: %d != %d",
			user.Purse, expected)
	}

	user.Change(val)

	s.Save(user)

	return nil
}
