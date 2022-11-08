package history

import (
	"database/sql"
	"time"

	"github.com/jbchouinard/multitool/config"
	"github.com/jbchouinard/multitool/database"
	"github.com/jbchouinard/multitool/errored"
)

var enabled bool

func init() {
	enabled = config.GetOptionDefault("history", "yes") == "yes"
	errored.Fatal("init db.history", database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec("CREATE TABLE IF NOT EXISTS history (ts TIMESTAMP, key TEXT, value TEXT)")
		return err
	}))
}

type Entry struct {
	Timestamp time.Time
	Key       string
	Value     string
}

func Add(k string, v string) {
	if !enabled {
		return
	}
	ts := time.Now().UTC()
	errored.Fatal("history.add", database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO history (ts, key, value) VALUES (?, ?, ?)",
			ts, k, v,
		)
		return err
	}))
}

func Purge(asOf time.Time) {
	errored.Fatal("history.purge", database.Tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"DELETE FROM history WHERE ts < ?",
			asOf,
		)
		return err
	}))
}

func GetLast(k string, n int) []*Entry {
	values := make([]*Entry, 0)

	errored.Fatal("history.getLast", database.Tx(func(tx *sql.Tx) error {
		rows, err := tx.Query(
			"SELECT ts, key, value FROM history WHERE key=? ORDER BY ts DESC LIMIT ?",
			k, n,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var ts time.Time
			var key string
			var value string
			if err := rows.Scan(&ts, &key, &value); err != nil {
				return err
			}
			values = append(values, &Entry{ts, key, value})
		}
		return rows.Err()
	}))
	return values
}