package bank

import "github.com/sitano/gobanktest/db"

type Bank interface {
	Balances() db.View

	// Notes
	// - Transaction value can be negative.
	// ­ All users starts with balance value set 100.
	// ­ Negative balance can’t be increased ­ only decreased to positive balance.
	Transaction(id db.UserId, val int64) error
}
