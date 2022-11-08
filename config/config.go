package config

import (
	"database/sql"

	"github.com/jbchouinard/multitool/database"
	"github.com/jbchouinard/multitool/errored"
)

func init() {
	errored.Fatal("init db.settings", database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)")
		return err
	}))
}

func GetOption(k string) string {
	var v string
	errored.Fatal("get option", database.Tx(func(tx *sql.Tx) error {
		err := tx.QueryRow("SELECT value FROM settings WHERE key=?", k).Scan(&v)
		if err != nil && err != sql.ErrNoRows {
			return err
		} else {
			return nil
		}
	}))
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
	errored.Fatal("set option", database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			`INSERT INTO settings (key, value) VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value`,
			key, value,
		)
		return err
	}))
}
