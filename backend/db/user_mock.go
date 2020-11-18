package db

import (
	"fmt"
	"sync"

	"github.com/enoperm/internet-services-2020/model"
)

type MockUserDB struct {
	users map[string]model.User
	mut   sync.RWMutex
}

func (db *MockUserDB) FetchUser(username string) (model.User, error) {
	db.mut.RLock()
	defer db.mut.RUnlock()

	value, ok := db.users[username]
	if !ok {
		return model.User{}, fmt.Errorf("could not find user by given name")
	}
	return value, nil
}

func (db *MockUserDB) InsertUser(user model.User) error {
	db.mut.Lock()
	defer db.mut.Unlock()

	if _, alreadyExists := db.users[user.Name]; alreadyExists {
		return model.ErrRegisterUsernameUnavailable
	}

	db.users[user.Name] = user
	return nil
}
