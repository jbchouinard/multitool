package database

import (
	"context"
	"database/sql"
)

func TxExec(query string, args ...any) (sql.Result, error) {
	conn := Conn()
	tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(query, args...)
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

type TxQueryFunc func(func(*sql.Rows) error) error

func TxQuery(query string, args ...any) TxQueryFunc {
	return func(forEach func(*sql.Rows) error) error {
		conn := Conn()
		tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback()

		rows, err := tx.Query(query, args...)
		if err != nil {
			return err
		}
		for rows.Next() {
			err := forEach(rows)
			if err != nil {
				return err
			}
		}
		err = rows.Err()
		if err == nil {
			return tx.Commit()
		} else {
			return err
		}
	}
}

type TxQueryRowFunc func(func(*sql.Row) error) error

func TxQueryRow(query string, args ...any) TxQueryRowFunc {
	return func(fRow func(*sql.Row) error) error {
		conn := Conn()
		tx, err := conn.BeginTx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return err
		}
		defer tx.Rollback()
		row := tx.QueryRow(query, args...)
		err = fRow(row)
		if err == nil {
			return tx.Commit()
		} else {
			return err
		}
	}
}
