package http

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jbchouinard/wmt/database"
	"github.com/jbchouinard/wmt/env"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/template"
)

func init() {
	createRequestTables()
}

type RequestDefinition struct {
	id       int64
	Name     string
	Method   string
	URL      string
	Query    map[string]string
	Headers  map[string]string
	Template *template.FileTemplate
	Body     []byte
}

func NewRequestDefinition(name string, method string, url string) *RequestDefinition {
	return &RequestDefinition{
		Name: name, Method: method, URL: url,
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}
}

func (rd *RequestDefinition) Eval(e *env.Env, ps []*env.KV) *RequestDefinition {
	eval := func(part string, value string) string {
		s, err := template.EvalString(fmt.Sprintf("request.%s.%s", rd.Name, part), value, e, ps)
		errored.Check(err, "")
		return s
	}
	evaled := NewRequestDefinition(rd.Name+".eval", eval("method", rd.Method), eval("url", rd.URL))
	if rd.Template != nil {
		var err error
		evaled.Body, err = rd.Template.Eval(e, ps)
		errored.Check(err, "")
	} else if rd.Body != nil {
		evaled.Body = rd.Body
	}
	i := 0
	for rawName, rawVal := range rd.Headers {
		name := eval(fmt.Sprintf("header[%d].name", i), rawName)
		val := eval(fmt.Sprintf("header[%d].value", i), rawVal)
		evaled.Headers[name] = val
		i++
	}
	i = 0
	for rawName, rawVal := range rd.Query {
		name := eval(fmt.Sprintf("query[%d].name", i), rawName)
		val := eval(fmt.Sprintf("query[%d].value", i), rawVal)
		evaled.Query[name] = val
		i++
	}
	return evaled
}

func (rd *RequestDefinition) Request() *http.Request {
	var body io.Reader
	if rd.Body != nil {
		body = bytes.NewBuffer(rd.Body)
	}
	request, err := http.NewRequest(rd.Method, rd.URL, body)
	errored.Check(err, "")
	for name, val := range rd.Headers {
		request.Header[name] = []string{val}
	}
	query := request.URL.Query()
	for name, val := range rd.Query {
		query.Add(name, val)
	}
	request.URL.RawQuery = query.Encode()
	return request
}

func (r *RequestDefinition) String() string {
	return fmt.Sprintf("%s %s %s", r.Name, r.Method, r.URL)
}

func (r *RequestDefinition) Details() string {
	lines := make([]string, 0, 4)
	lines = append(lines, fmt.Sprintf("%s %s", r.Method, r.URL), "  query:")
	for k, v := range r.Query {
		lines = append(lines, fmt.Sprintf("    %s: %s", k, v))
	}
	lines = append(lines, "  headers:")
	for k, v := range r.Headers {
		lines = append(lines, fmt.Sprintf("    %s: %s", k, v))
	}
	lines = append(lines, fmt.Sprintf("  body: %s", r.Template))
	return strings.Join(lines, "\n")
}

func (r *RequestDefinition) templateID() *int64 {
	if r.Template == nil {
		return nil
	} else {
		id := r.Template.Id()
		return &id
	}
}

func (r *RequestDefinition) Save() error {
	return database.Tx(func(tx *sql.Tx) error {
		query := `
		INSERT INTO request_definition(name, method, url, template_id)
			VALUES(?, ?, ?, ?)
			ON CONFLICT(name) DO UPDATE SET
				method=excluded.method,
				url=excluded.url,
				template_id=excluded.template_id
		`
		res, err := tx.Exec(query, r.Name, r.Method, r.URL, r.templateID())
		if err != nil {
			return err
		}
		insertId, err := res.LastInsertId()
		if insertId != 0 {
			r.id = insertId
		}
		if err != nil {
			return err
		}
		if err := saveMap(tx, "request_header", r.id, r.Headers); err != nil {
			return err
		}
		if err := saveMap(tx, "request_query", r.id, r.Query); err != nil {
			return err
		}
		return nil
	})
}

func (r *RequestDefinition) Delete() error {
	_, err := database.TxExec("DELETE FROM request_definition WHERE id=?", r.id)
	return err
}

func ListRequestDefinitions() []string {
	names := make([]string, 0)
	err := database.TxQuery("SELECT name FROM request_definition")(func(row *sql.Rows) error {
		var name string
		if err := row.Scan(&name); err != nil {
			return err
		}
		names = append(names, name)
		return nil
	})
	errored.Check(err, "")
	return names
}

func LoadRequestDefinition(name string) (*RequestDefinition, error) {
	req := &RequestDefinition{0, name, "", "", nil, nil, nil, nil}
	var templateName *string
	err := database.Tx(func(tx *sql.Tx) error {
		err := tx.QueryRow(`
		SELECT r.id, r.method, r.url, t.name
			FROM request_definition r
			LEFT OUTER JOIN template t
				ON r.template_id = t.id
			WHERE r.name = ?`,
			req.Name,
		).Scan(&req.id, &req.Method, &req.URL, &templateName)
		if err != nil {
			return err
		}
		req.Headers, err = loadMap(tx, "request_header", req.id)
		if err != nil {
			return err
		}
		req.Query, err = loadMap(tx, "request_query", req.id)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	if templateName != nil {
		req.Template, err = template.SelectTemplate(*templateName)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func loadMap(tx *sql.Tx, table string, reqid int64) (map[string]string, error) {
	m := make(map[string]string)
	query := fmt.Sprintf("SELECT name, value FROM %s WHERE request_definition_id=?", table)
	rows, err := tx.Query(query, reqid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		var value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		m[name] = value
	}
	return m, rows.Err()
}

func saveMap(tx *sql.Tx, table string, reqid int64, m map[string]string) error {
	_, err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE request_definition_id=?", table), reqid)
	if err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}
	values := make([]string, 0, len(m))
	args := make([]any, 0, 3*len(m))
	for name, val := range m {
		values = append(values, "(?, ?, ?)")
		args = append(args, reqid, name, val)
	}
	query := fmt.Sprintf(
		"INSERT INTO %s(request_definition_id, name, value) VALUES%s",
		table, strings.Join(values, ","),
	)
	_, err = tx.Exec(query, args...)
	return err
}

func createRequestTables() {
	_, err := database.TxExec(
		`CREATE TABLE IF NOT EXISTS request_definition (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			method TEXT,
			url TEXT,
			template_id INTEGER,
			FOREIGN KEY(template_id) REFERENCES template(id)
		)`,
	)
	errored.Check(err, "init db.request_definition")

	mapTableQuery := `CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY,
		request_definition_id INTEGER,
		name TEXT,
		value TEXT,
		FOREIGN KEY(request_definition_id) REFERENCES request_definition(id) ON DELETE CASCADE
	)`
	_, err = database.TxExec(fmt.Sprintf(mapTableQuery, "request_header"))
	errored.Check(err, "init db.request_header")
	_, err = database.TxExec(fmt.Sprintf(mapTableQuery, "request_query"))
	errored.Check(err, "init db.request_query")
}
