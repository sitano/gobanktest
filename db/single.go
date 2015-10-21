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

func (s *inMemory) Load(id UserId) (*User, error) {
	purse, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("There is no such user %d", id)
	}

	return &User{
		Id: id,
		Purse: purse,
	}, nil
}

func (s *inMemory) Save(user *User) error {
	s.data[user.Id] = user.Purse
	return nil
}

func (s *inMemory) List() View {
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
	if _, ok := s.data[user.Id]; ok {
		return fmt.Errorf("Can't put user %d into the storage: already present", user.Id)
	}

	return s.Save(user)
}

// Single threaded implementation has atomicity guarantee by definition
func (s *inMemory) Change(id UserId, val int64, expected Purse) error {
	user, err := s.Load(id)
	if err != nil {
		return err
	}

	if user.Purse != expected {
		return fmt.Errorf("User balance have been changed since. %v != %d",
			user, expected)
	}

	user.Change(val)

	s.Save(user)

	return nil
}
