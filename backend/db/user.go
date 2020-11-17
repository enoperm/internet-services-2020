package db

import (
    "database/sql"
    "github.com/enoperm/internet-services-2020/model"

	sqlite "github.com/mattn/go-sqlite3"
)

// TODO: Interface is missing Context versions of the API
type UserDatabase interface {
    InitializeUserSchema() error

    FetchUser(username string) (model.User, error)
    InsertUser(model.User) error
}
var _ UserDatabase = &ApplicationDatabase{}

func (db ApplicationDatabase) InitializeUserSchema() error {
    _, err := db.Db.Exec(`
        CREATE TABLE
        IF NOT EXISTS
        users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(48) UNIQUE,
            pwsalt BLOB,
            pwhash BLOB
        );
    `)
    logger.Println("dbg:", err)
    return err
}


func (db ApplicationDatabase) FetchUser(username string) (model.User, error) {
    row := db.Db.QueryRow(`SELECT id, name, pwsalt, pwhash FROM users WHERE name = ?`, username)
    user := model.User{
        Password: &model.Password{},
    }

	var salt, hash []byte

	err := row.Scan(&user.Id, &user.Name, &salt, &hash)
    switch {
    case err == sql.ErrNoRows:
        return user, err

    case err != nil:
        logger.Printf("users/fetch: %s: %s", username, err)
        return user, err
    }

	copy(user.Password.Salt[:], salt)
	copy(user.Password.Hash[:], hash)
    return user, nil
}

func (db ApplicationDatabase) InsertUser(user model.User) error {
    _, err := db.Db.Exec(`
        INSERT OR FAIL
        INTO users (name, pwsalt, pwhash)
        VALUES (?, ?, ?)
    `, user.Name, user.Password.Salt[:], user.Password.Hash[:])
    switch err.(type) {
	case sqlite.Error:
		if err.(sqlite.Error).Code == sqlite.ErrConstraint {
			return model.ErrRegisterUsernameUnavailable
		}
    default:
    }

    return err
}
