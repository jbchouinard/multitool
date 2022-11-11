package template

import (
	"bytes"
	"database/sql"
	"fmt"
	htmltemplate "html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jbchouinard/wmt/database"
	"github.com/jbchouinard/wmt/env"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/path"
)

var templateDir string

func init() {
	createTemplateTable()
	templateDir = filepath.Join(path.WorkDir, "templates")
	err := os.MkdirAll(templateDir, os.ModePerm)
	if err != nil {
		errored.Fatalf("failed to create config dir %q: %s", templateDir, err)
	}
}

func TemplatePath(name string) string {
	return filepath.Join(templateDir, name+".tmpl")
}

type FileTemplate struct {
	id           int64
	Name         string
	IsHTML       bool
	HTMLTemplate *htmltemplate.Template
	TextTemplate *template.Template
}

func (t *FileTemplate) Id() int64 {
	return t.id
}

func (t *FileTemplate) Path() string {
	return TemplatePath(t.Name)
}

func (t *FileTemplate) String() string {
	if t == nil {
		return "empty"
	}
	var kind string
	if t.IsHTML {
		kind = "html"
	} else {
		kind = "text"
	}
	return fmt.Sprintf("%s %s@%s", t.Name, kind, t.Path())
}

func newTemplate(id int64, name string, isHtml bool) *FileTemplate {
	return &FileTemplate{id, name, isHtml, nil, nil}
}

func (t *FileTemplate) Parse() error {
	contentBytes, err := ioutil.ReadFile(TemplatePath(t.Name))
	if err != nil {
		return err
	}
	content := string(contentBytes)
	if t.IsHTML {
		tmpl, err := htmltemplate.New(t.Name).Option("missingkey=error").Parse(content)
		if err != nil {
			return err
		}
		t.HTMLTemplate = tmpl
	} else {
		tmpl, err := template.New(t.Name).Option("missingkey=error").Parse(content)
		if err != nil {
			return err
		}
		t.TextTemplate = tmpl
	}
	return nil
}

func EvalString(name string, tmpl string, e *env.Env, ps []*env.KV) (string, error) {
	args := MakeArgs(e, ps)
	if strings.Contains(tmpl, "{{") {
		sTemplate, err := template.New(name).Option("missingkey=error").Parse(tmpl)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := sTemplate.Execute(&buf, args); err != nil {
			return "", err
		}
		return string(buf.Bytes()), nil
	} else {
		return tmpl, nil
	}
}

func MakeArgs(e *env.Env, ps []*env.KV) map[string]string {
	args := e.GetAll()
	env.AddKVs(args, ps)
	return args
}

func (t *FileTemplate) Eval(e *env.Env, params []*env.KV) ([]byte, error) {
	args := MakeArgs(e, params)
	if t.IsHTML {
		if t.HTMLTemplate == nil {
			if err := t.Parse(); err != nil {
				return nil, err
			}
		}
		buf := new(bytes.Buffer)
		if err := t.HTMLTemplate.Execute(buf, args); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	} else {
		if t.TextTemplate == nil {
			if err := t.Parse(); err != nil {
				return nil, err
			}
		}
		buf := new(bytes.Buffer)
		if err := t.TextTemplate.Execute(buf, args); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func createTemplateTable() {
	_, err := database.TxExec(
		`CREATE TABLE IF NOT EXISTS template (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE,
			is_html INTEGER
		)`,
	)
	errored.Check(err, "init db.template")
}

func CreateTemplate(name string, isHtml bool) (*FileTemplate, error) {
	tmpl := newTemplate(0, name, isHtml)
	res, err := database.TxExec(
		`INSERT INTO template(name, is_html)
			VALUES(?, ?)`,
		name, isHtml,
	)
	if err != nil {
		return nil, err
	}
	if id, err := res.LastInsertId(); err != nil {
		return nil, err
	} else {
		tmpl.id = id
	}
	return tmpl, nil
}

func (t *FileTemplate) Update() error {
	_, err := database.TxExec(
		`UPDATE template SET
			name=?, is_html=?
			WHERE id=?`,
		t.Name, t.IsHTML, t.id,
	)
	return err
}

func (t *FileTemplate) Delete() {
	_, err := database.TxExec("DELETE FROM template WHERE id=?", t.id)
	errored.Check(err, "delete template")
}

func DeleteTemplate(name string) {
	_, err := database.TxExec("DELETE FROM template WHERE name=?", name)
	errored.Check(err, "delete template")
}

func SelectTemplate(name string) (*FileTemplate, error) {
	var id int64
	var isHtml bool

	err := database.TxQueryRow(
		`SELECT id, is_html FROM template WHERE name=?`,
		name,
	)(func(row *sql.Row) error {
		return row.Scan(&id, &isHtml)
	})
	if err != nil {
		return nil, err
	}
	return newTemplate(id, name, isHtml), nil
}

func ListTemplates() []*FileTemplate {
	templates := make([]*FileTemplate, 0)
	err := database.TxQuery("SELECT id, name, is_html FROM template")(func(row *sql.Rows) error {
		var id int64
		var name string
		var isHtml bool
		err := row.Scan(&id, &name, &isHtml)
		if err != nil {
			return err
		}
		templates = append(templates, newTemplate(id, name, isHtml))
		return nil
	})
	errored.Check(err, "list templates")
	return templates
}
