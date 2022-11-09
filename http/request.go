package http

import (
	"bytes"
	"database/sql"
	htmltemplate "html/template"
	"text/template"

	"github.com/jbchouinard/wmt/database"
	"github.com/jbchouinard/wmt/errored"

	// template table must be created before request
	_ "github.com/jbchouinard/wmt/template"
)

func init() {
	createRequestTables()
}

type Body interface {
	Bytes() ([]byte, error)
}

type Environment map[string]string

type TemplatedBody struct {
	tmpl *template.Template
	env  Environment
}

func (b *TemplatedBody) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := b.tmpl.Execute(buf, b.env); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type HTMLTemplatedBody struct {
	tmpl *htmltemplate.Template
	env  Environment
}

func (b *HTMLTemplatedBody) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := b.tmpl.Execute(buf, b.env); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type RequestTemplate struct {
	id      int
	Verb    string
	Path    string
	Query   map[string]string
	Headers map[string]string
	Body    Body
}

func LoadRequestTemplate(name string) (*RequestTemplate, error) {
	tmpl := new(RequestTemplate)
	err := database.Tx(func(tx *sql.Tx) error {
		tx.QueryRow(`SELECT id, name, verb, path, body`)
		return nil
	})
	return tmpl, err
}

func createRequestTables() {
	_, err := database.TxExec(
		`CREATE TABLE IF NOT EXISTS request_template (
			id INTEGER PRIMARY KEY
			name TEXT UNIQUE,
			verb TEXT,
			path TEXT,
			body INTEGER,
			FOREIGN KEY(body) REFERENCES template(id)
		)`,
	)
	errored.Check(err, "init db.request_template")
	_, err = database.TxExec(
		`CREATE TABLE IF NOT EXISTS request_header (
			id INTEGER PRIMARY KEY
			request_template_id INTEGER,
			name TEXT,
			value TEXT,
			FOREIGN KEY(request_template_id) REFERENCES request_template(id) ON DELETE CASCADE
		)`,
	)
	errored.Check(err, "init db.request_header")
	_, err = database.TxExec(
		`CREATE TABLE IF NOT EXISTS request_header (
			id INTEGER PRIMARY KEY
			request_id TEXT UNIQUE,
			name TEXT,
			value TEXT,
			FOREIGN KEY(request_template_id) REFERENCES request_template(id) ON DELETE CASCADE
		)`,
	)
	errored.Check(err, "init db.request_query")
}
