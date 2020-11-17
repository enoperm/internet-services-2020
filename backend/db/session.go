package db

import (
	"time"
    "database/sql"
    "github.com/enoperm/internet-services-2020/model"
)

// TODO: Interface is missing Context versions of the API
type SessionDatabase interface {
    InitializeSessionSchema() error

	InsertSession(sess model.Session) (*model.Session, error)
	FetchSession(sid int64) (*model.Session, error)
	RemoveSession(sid int64) error

	TouchSession(sid int64) error
}
var _ SessionDatabase = &ApplicationDatabase{}

func (db ApplicationDatabase) InitializeSessionSchema() error {
    _, err := db.Db.Exec(`
        CREATE TABLE
        IF NOT EXISTS
        sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
        	hmac BLOB,
			created_at STRING,
			last_seen STRING,

			FOREIGN KEY(user_id) REFERENCES users (id)
        );
    `)
	logger.Println("sessions/init-schema:", err)
    return err
}


func (db ApplicationDatabase) FetchSession(sid int64) (*model.Session, error) {
	row := db.Db.QueryRow(`SELECT user_id, hmac, created_at, last_seen FROM sessions WHERE id = ?`, sid)
    sess := model.Session{
		SessionID: sid,
    }
	var cat, ls string
    err := row.Scan(&sess.UserID, &sess.Mac, &cat, &ls)

    switch {
    case err == sql.ErrNoRows:
        return nil, err

    case err != nil:
        logger.Printf("sessions/fetch: sid(%v): %s", sid, err)
        return nil, err
    }

	sess.CreatedAt, err = time.Parse(ISO8601, cat)
	if err != nil { return nil, err }

	sess.LastSeen, err = time.Parse(ISO8601, ls)
	if err != nil { return nil, err }

    return &sess, nil
}

func (db ApplicationDatabase) InsertSession(sess model.Session) (*model.Session, error) {
    res, err := db.Db.Exec(`
        INSERT OR FAIL
        INTO sessions (user_id, hmac, created_at, last_seen)
        VALUES (?, ?, datetime(), datetime())
    `, sess.UserID, sess.Mac)

	if err == nil {
		sess.SessionID, _ = res.LastInsertId()
	}

    return &sess, err
}

func (db ApplicationDatabase) RemoveSession(sid int64) error {
	return nil
}

// TODO: expose this function to session middleware.
func (db ApplicationDatabase) TouchSession(sid int64) error {
	_, err := db.Db.Exec(`
		UPDATE sessions
		SET last_seen = datetime()
		WHERE id = ?
	`, sid)
	return err
}
