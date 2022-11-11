package config

import (
	"database/sql"
	"fmt"

	"github.com/jbchouinard/wmt/database"
	"github.com/jbchouinard/wmt/errored"
)

func init() {
	_, err := database.TxExec("CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)")
	errored.Check(err, "init db.settings")
}

// These values are set by packages that use the config system
var DefaultValues map[string]string = map[string]string{}
var ValidValues map[string]map[string]bool = map[string]map[string]bool{}

type Option struct {
	Key   string
	Value string
}

func GetAll() []*Option {
	options := make([]*Option, 0, len(DefaultValues))
	for key := range DefaultValues {
		options = append(options, &Option{key, Get(key)})
	}
	return options
}

func Get(k string) string {
	var v string
	err := database.TxQueryRow("SELECT value FROM settings WHERE key=?", k)(func(row *sql.Row) error {
		err := row.Scan(&v)
		if err == sql.ErrNoRows {
			v = DefaultValues[k]
			return nil
		}
		return err
	})
	errored.Check(err, "get option")
	return v
}

func Set(key string, value string) error {
	if DefaultValues[key] == "" {
		return fmt.Errorf("invalid option key %q", key)
	}
	valid := ValidValues[key]
	if valid != nil {
		if !valid[value] {
			return fmt.Errorf("invalid option value %q", value)
		}
	}
	_, err := database.TxExec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
			ON CONFLICT(key) DO UPDATE SET value=excluded.value`,
		key, value,
	)
	errored.Check(err, "set option")
	return nil
}

func Unset(key string) error {
	if DefaultValues[key] == "" {
		return fmt.Errorf("invalid option key %q", key)
	}
	_, err := database.TxExec("DELETE FROM settings WHERE key=?", key)
	errored.Check(err, "unset option")
	return nil
}
