package db

import (
    "database/sql"
    "log"
    "os"
)

var (
	logger *log.Logger
	_ = logger
)

func init() {
	logger = log.New(os.Stdout, "db: ", log.LstdFlags | log.LUTC | log.Lmsgprefix)
}

// This is hardly idiomatic, but it'll do for now.
type ApplicationDatabase struct {
    Db *sql.DB
}
