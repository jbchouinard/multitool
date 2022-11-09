package database

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/jbchouinard/multitool/errored"
	"github.com/jbchouinard/multitool/path"
	_ "github.com/mattn/go-sqlite3"
)

var dbFile string
var dbConn *sql.Conn

func Conn() *sql.Conn {
	if dbConn == nil {
		db, err := sql.Open("sqlite3", dbFile)
		errored.Check(err, "db open")
		db.SetConnMaxIdleTime(0)
		db.SetConnMaxLifetime(0)
		db.SetMaxOpenConns(1)
		db.SetMaxOpenConns(1)
		dbConn, err = db.Conn(context.Background())
		errored.Check(err, "db connect")
	}
	return dbConn
}

func Tx(f func(tx *sql.Tx) error) error {
	conn := Conn()
	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = f(tx)
	if err == nil {
		return tx.Commit()
	} else {
		return err
	}
}

func init() {
	dbFile = filepath.Join(path.WorkDir, "multitool.db")
}
