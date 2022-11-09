package config

import (
	"database/sql"
	"fmt"

	"github.com/jbchouinard/multitool/database"
	"github.com/jbchouinard/multitool/errored"
)

func init() {
	err := database.TxExec("CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)")
	errored.Check(err, "init db.settings")
}

var defaultValues map[string]string = map[string]string{
	"clipboard": "no",
	"history":   "yes",
	"editor":    "nano",
}

var validValues map[string]map[string]bool = map[string]map[string]bool{
	"clipboard": {"yes": true, "no": true},
	"history":   {"yes": true, "no": true},
}

type Option struct {
	Key   string
	Value string
}

func GetAll() []*Option {
	options := make([]*Option, 0, len(defaultValues))
	for key := range defaultValues {
		options = append(options, &Option{key, Get(key)})
	}
	return options
}

func Get(k string) string {
	var v string
	err := database.TxQueryRow("SELECT value FROM settings WHERE key=?", k)(func(row *sql.Row) error {
		err := row.Scan(&v)
		if err == sql.ErrNoRows {
			v = defaultValues[k]
			return nil
		}
		return err
	})
	errored.Check(err, "get option")
	return v
}

func Set(key string, value string) error {
	if defaultValues[key] == "" {
		return fmt.Errorf("invalid option key %q\n", key)
	}
	valid := validValues[key]
	if valid != nil {
		if !valid[value] {
			return fmt.Errorf("invalid option value %q\n", value)
		}
	}
	err := database.TxExec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value`,
		key, value,
	)
	errored.Check(err, "set option")
	return nil
}

func Unset(key string) error {
	if defaultValues[key] == "" {
		return fmt.Errorf("invalid option key %q\n", key)
	}
	err := database.TxExec("DELETE FROM settings WHERE key=?", key)
	errored.Check(err, "unset option")
	return nil
}
