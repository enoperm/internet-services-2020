package db

import (
	"database/sql"
	"github.com/enoperm/internet-services-2020/model"

	sqlite "github.com/mattn/go-sqlite3"
)

// TODO: Interface is missing Context versions of the API
type ProfileDatabase interface {
	InitializeProfileSchema() error

	FetchProfile(uid int64) (model.Profile, error)
	UpdateProfile(uid int64, newProfile model.Profile) error
}

var _ ProfileDatabase = &ApplicationDatabase{}

func (db ApplicationDatabase) InitializeProfileSchema() error {
	_, err := db.Db.Exec(`
        CREATE TABLE
        IF NOT EXISTS
        profiles (
            user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			last_smoke STRING,
			daily_average INTEGER,
			sticks_per_pack INTEGER,
			price_per_pack INTEGER,
			start_year INTEGER,

			FOREIGN KEY(user_id) REFERENCES users (id)
        );
    `)
	logger.Println("dbg:", err)
	return err
}

func (db ApplicationDatabase) FetchProfile(uid int64) (model.Profile, error) {
	row := db.Db.QueryRow(
		`SELECT last_smoke,
				daily_average,
			    sticks_per_pack,
			    price_per_pack,
			    start_year
		 FROM profiles
		 WHERE user_id = ?`, uid)
	var profile model.Profile

	err := row.Scan(&profile.LastSmoke,
				    &profile.DailyAverage,
				    &profile.SticksPerPack,
				    &profile.PricePerPack,
				    &profile.StartYear)
	switch {
	case err == sql.ErrNoRows:
		return profile, err

	case err != nil:
		logger.Printf("profiles/fetch: %s: %s", uid, err)
		return profile, err
	}

	return profile, nil
}

func (db ApplicationDatabase) UpdateProfile(uid int64, profile model.Profile) error {
	logger.Println(profile)
	_, err := db.Db.Exec(`
        INSERT OR REPLACE
        INTO profiles (user_id, last_smoke, daily_average, sticks_per_pack, price_per_pack, start_year)
        VALUES (?, ?, ?, ?, ?, ?)
    `, uid, profile.LastSmoke, profile.DailyAverage, profile.SticksPerPack, profile.PricePerPack, profile.StartYear)
	switch err.(type) {
	case sqlite.Error:
		if err.(sqlite.Error).Code == sqlite.ErrConstraint {
			return model.ErrRegisterUsernameUnavailable
		}
	default:
	}

	return err
}