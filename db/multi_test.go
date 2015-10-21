package db

import (
	"testing"
	"runtime"
	"sync"
)

var threads int = 8
var steps int = 100

const max_iters = 100

func simple_inc(db Storage, name UserId) error {
	user, err := db.Load(name)
	if err != nil {
		return err
	}

	user.Purse ++

	return db.Save(user)
}

func simple_change_by1(db Storage, name UserId) error {
	user, err := db.Load(name)
	if err != nil {
		return err
	}

	return db.Transaction().Change(user.Id, 1, user.Purse)
}

func spin_change_by1(db Storage, name UserId) error {
	for {
		user, err := db.Load(name)
		if err != nil {
			return err
		}

		if err := db.Transaction().Change(user.Id, 1, user.Purse); err == nil {
			break
		}
	}

	return nil
}

// That is far from perfect, but works for now
func run_mt_test(
		t *testing.T, db Storage,
		threads int, steps int, u *User,
		finished func (u *User, iter int) bool,
		step func(db Storage, step int) error) {

	runtime.GOMAXPROCS(threads)

	if threads < 2 || runtime.NumCPU() < 2 {
		t.Fatal("This test requires more than 1 physical core to be present")
	}

	wait := sync.WaitGroup{}
	iteration := 0

	if err := db.Save(u); err != nil {
		t.Fatal("DB failed to create a user", u, err)
	}

	for {
		iteration ++

		t.Log("Go for iteration", iteration)

		wait.Add(threads)

		for i := 0; i < threads; i ++ {
			go func() {
				defer wait.Done()

				for j := 0; j < iteration * steps; j ++ {
					if err := step(db, 1 + j); err != nil {
						t.Log("Error occured in", i, "th stepper:", err)
						return
					}
				}
			}()
		}

		wait.Wait()

		u1, err := db.Load(u.Id)
		if err != nil {
			t.Fatal("Load failed for", u.Id, err)
		}

		if finished(u1, iteration) {
			t.Log("User now have", u1.Purse)
			break
		}

		if iteration >= max_iters {
			t.Fatal("Max count of iterations reached")
		}

		if err = db.Save(u); err != nil {
			t.Fatal("Reset failed", u, err)
		}
	}
}

func Test_STSaveShouldNotWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse != Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return simple_inc(db, u.Id)
		})
}

func Test_STSimpleChangeShouldNotWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse != Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return simple_change_by1(db, u.Id)
		})
}

func Test_STSpinnedChangeShouldNotWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse != Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return spin_change_by1(db, u.Id)
		})
}

func Test_MTSaveShouldNotWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryMTSafeStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse != Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return simple_inc(db, u.Id)
		})
}

func Test_MTSimpleChangeShouldNotWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryMTSafeStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse != Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return simple_change_by1(db, u.Id)
		})
}

func Test_MTSpinnedChangeShouldWorkInMTEnv(t *testing.T) {
	u := &User{Id: 1, Purse: 0}

	run_mt_test(t, NewInMemoryMTSafeStorage(), threads, steps, u,
		func(u *User, iter int) bool {
			return u.Purse > 0 && u.Purse == Purse(iter * threads * steps)
		},
		func (db Storage, i int) error {
			return spin_change_by1(db, u.Id)
		})
}

func Test_MTEnvDeadlocksTest(t *testing.T) {
	// TODO
}
