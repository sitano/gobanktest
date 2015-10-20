package db

import "strconv"

type UserId string
type Purse int64

type User struct {
	Name UserId
	Purse Purse
}

type Storage interface {
	Load(name UserId) (*User, error)

	// Just to show this would not work for our scenario
	Save(user *User) error

	// Simple transaction abstraction hinting we need ACID guarantees
	Transaction() Tx
}

type Tx interface {
	// Put a User into the DB if it is not present
	PutIfAbsent(user *User) error

	// Change purse balance of the specified User if
	// it has verified balance (the balance did not change
	// since last read).
	//
	// It is possible to put this peace of business
	// logic here into kind of storage procedure, but
	// i choosed not to go this way.
	Change(name UserId, val int64, expected Purse) error
}

func Compare(u1, u2 *User) bool {
	return u1 == u2 || (u1.Name == u2.Name && u1.Purse == u2.Purse)
}

func (u *User) String() string {
	return "User{" + string(u.Name) + ", " + strconv.FormatInt(int64(u.Purse), 10) + "}"
}