package db

import "strconv"

type UserId int64
type Purse int64
type View map[UserId]Purse

type User struct {
	Id UserId
	Purse Purse
}

type Storage interface {
	Load(id UserId) (*User, error)

	// Just to show this would not work for our scenario
	Save(user *User) error

	List() View

	// Separate simple transaction abstraction.
	// Need ACID guarantees to perform transactional changes to Purses.
	Transaction() Tx
}

type Tx interface {
	// Put a User into the DB if it is not present
	PutIfAbsent(user *User) error

	// Change purse balance of the specified User if
	// it has verified balance (the balance did not change
	// since last read).
	//
	// Available options here:
	// - put whole transaction logic into db
	// - make start tx/commit acid with locks
	// - single cas operation
	//
	// I provide CAS op like here to show how to
	// organize purses changes using that simple
	// basic op.
	Change(id UserId, val int64, expected Purse) error
}

func Compare(u1, u2 *User) bool {
	return u1 == u2 || (u1.Id == u2.Id && u1.Purse == u2.Purse)
}

func (u *User) Change(val int64) {
	u.Purse = Purse(int64(u.Purse) + val)
}

func (u *User) String() string {
	return "[User " +
		strconv.FormatInt(int64(u.Id), 10) + ": " +
		strconv.FormatInt(int64(u.Purse), 10) + "]"
}