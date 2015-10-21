package db

import "testing"

func Test_CompareUsers(t *testing.T) {
	u1 := &User{Id: 1, Purse: 1}
	u2 := &User{Id: 1, Purse: 2}

	if Compare(u1, u2) {
		t.Error("Users must not be equal", u1, u2)
	}

	u2.Purse = 1
	if !Compare(u1, u2) {
		t.Error("Users must be equal", u1, u2)
	}

	u2.Id = 2
	if Compare(u1, u2) {
		t.Error("Users must not be equal", u1, u2)
	}
}

func Test_STInit(t *testing.T) {
	if NewInMemoryStorage() == nil {
		t.Error("Storage failed to init")
	}
}

func Test_STLoadOfMissingUserShouldFail(t *testing.T) {
	db := NewInMemoryStorage()

	if user, err := db.Load(1); user != nil {
		t.Fatal("DB should not load non existent user")
	} else if err.Error() != "There is no such user 1" {
		t.Error("DB returned some unknown error", err)
	}
}

func Test_STSaveShouldCreateUser(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 1}

	if err := db.Save(u1); err != nil {
		t.Fatal("DB failed to create a user", u1, err)
	}
}

func Test_STSaveShouldRewriteExistingUser(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 1}
	u2 := &User{Id: 2, Purse: 3}
	u3 := &User{Id: 1, Purse: 2}

	if err := db.Save(u1); err != nil {
		t.Fatal("DB failed to create a user", u1, err)
	}

	if err := db.Save(u2); err != nil {
		t.Fatal("DB failed to create a user", u2, err)
	}

	if err := db.Save(u3); err != nil {
		t.Error("DB failed to rewrite existing user with", u3, err)
	}

	if u4, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1.Id, err)
	} else if !Compare(u4, u3) {
		t.Error("The loaded user is different to saved user", u4, u3)
	}

	// We can't check db size here, as we don't have this api
	// Just check other users did not change
	if u5, err := db.Load(u2.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u2.Id, err)
	} else if !Compare(u5, u2) {
		t.Error("The loaded user is different to saved user", u5, u2)
	}
}

func Test_STLoadShouldSeeExistedUser(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 1}

	if err := db.Save(u1); err != nil {
		t.Fatal("DB failed to create a user", u1, err)
	}

	if u11, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1, err)
	} else if !Compare(u11, u1) {
		t.Error("The loaded user is different to saved user", u11, u1)
	}
}

func Test_STTXPutShouldNotRewrite(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 1}
	u2 := &User{Id: 1, Purse: 2}

	if err := db.Transaction().PutIfAbsent(u1); err != nil {
		t.Fatal("DB tx failed to create a user", u1, err)
	}

	if err := db.Transaction().PutIfAbsent(u1); err == nil {
		t.Error("DB tx put should not rewrite existed user", u1, err)
	}

	if err := db.Transaction().PutIfAbsent(u2); err == nil {
		t.Error("DB tx put should not rewrite existed user", u1, err)
	}

	if u3, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1.Id, err)
	} else if !Compare(u3, u1) {
		t.Error("User must not be modified by db tx put", u3, u1)
	}
}

func Test_STTXChangeShouldNotSaveUnexpectedValue(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 1}

	if err := db.Transaction().PutIfAbsent(u1); err != nil {
		t.Fatal("DB tx failed to create a user", u1, err)
	}

	if err := db.Transaction().Change(u1.Id, 1, 2); err == nil {
		t.Error("DB tx change must not change user with unexpected value", u1, err)
	}

	if u3, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1.Id, err)
	} else if !Compare(u3, u1) {
		t.Error("User must not be modified by failed db tx change", u3, u1)
	}

	if err := db.Transaction().Change(u1.Id, 1, 1); err != nil {
		t.Error("DB tx change must change user with expected value", u1, err)
	}

	if u3, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1.Id, err)
	} else if u3.Purse != 2 {
		t.Error("User must have purse balance equals 2", u3, u1)
	}
}

func Test_STTXChangeShouldAllowNegative(t *testing.T) {
	db := NewInMemoryStorage()

	u1 := &User{Id: 1, Purse: 0}

	if err := db.Transaction().PutIfAbsent(u1); err != nil {
		t.Fatal("DB tx failed to create a user", u1, err)
	}

	if err := db.Transaction().Change(u1.Id, -1, 0); err != nil {
		t.Error("DB tx change must not change user with unexpected value", u1, err)
	}

	if err := db.Transaction().Change(u1.Id, -1, -1); err != nil {
		t.Error("DB tx change must not change user with unexpected value", u1, err)
	}

	if u3, err := db.Load(u1.Id); err != nil {
		t.Fatal("DB failed to load an existent user", u1.Id, err)
	} else if u3.Purse != -2 {
		t.Error("User must have purse balance equals -2", u3, u1)
	}
}
