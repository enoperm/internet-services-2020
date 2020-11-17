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

const ISO8601 = "2006-01-02 15:04:05"

func init() {
	logger = log.New(os.Stdout, "db: ", log.LstdFlags | log.LUTC | log.Lmsgprefix)
}

// This is hardly idiomatic, but it'll do for now.
type ApplicationDatabase struct {
    Db *sql.DB
}
