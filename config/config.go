package config

import (
	"database/sql"

	"github.com/jbchouinard/multitool/database"
	"github.com/jbchouinard/multitool/errored"
)

func init() {
	errored.Check(database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)")
		return err
	}), "init db.settings")
}

func GetOption(k string) string {
	var v string
	errored.Check(database.Tx(func(tx *sql.Tx) error {
		err := tx.QueryRow("SELECT value FROM settings WHERE key=?", k).Scan(&v)
		if err != nil && err != sql.ErrNoRows {
			return err
		} else {
			return nil
		}
	}), "get option")
	return v
}

func GetOptionDefault(k string, d string) string {
	v := GetOption(k)
	if v != "" {
		return v
	} else {
		return d
	}
}

func SetOption(key string, value string) {
	errored.Check(database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO settings (key, value) VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value`,
			key, value,
		)
		return err
	}), "set option")
}
