package env

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jbchouinard/wmt/config"
	"github.com/jbchouinard/wmt/database"
	"github.com/jbchouinard/wmt/errored"
)

func init() {
	createEnvTable()
	Global.Load()
	envs["global"] = Global
	config.DefaultValues["env"] = "global"
	Current = GetEnv(config.Get("env"))
}

type Env struct {
	name   string
	parent *Env
	values map[string]string
}

var Global *Env = &Env{"global", nil, make(map[string]string)}
var Current *Env

func (e *Env) Name() string {
	return e.name
}

func SetCurrentEnv(name string) {
	config.Set("env", name)
	Current = GetEnv(name)
}

func newLocalEnv(name string) *Env {
	return &Env{name, Global, make(map[string]string)}
}

var envs map[string]*Env = make(map[string]*Env, 1)

func GetEnv(name string) *Env {
	env, ok := envs[name]
	if ok {
		return env
	} else {
		env = newLocalEnv(name)
		env.Load()
		envs[name] = env
		return env
	}
}

func (e *Env) Get(key string) (val string, ok bool) {
	v, ok := e.values[key]
	if !ok && e.parent != nil {
		return e.parent.Get(key)
	} else {
		return v, ok
	}
}

func (e *Env) GetAll() map[string]string {
	if e.parent == nil {
		return e.values
	} else {
		values := e.parent.GetAll()
		for k, v := range e.values {
			values[k] = v
		}
		return values
	}
}

type EnvVar struct {
	Source string
	Key    string
	Value  string
}

func (v *EnvVar) String() string {
	return fmt.Sprintf("%-16s %-16s %s", v.Source, v.Key, v.Value)
}

func (e *Env) List() []*EnvVar {
	var vars []*EnvVar
	if e.parent != nil {
		vars = e.parent.List()
	} else {
		vars = make([]*EnvVar, 0, len(e.values))
	}
	for k, v := range e.values {
		vars = append(vars, &EnvVar{e.name, k, v})
	}
	return vars
}

func (e *Env) Set(key string, val string) {
	e.values[key] = val
}

func (e *Env) Unset(key string) {
	delete(e.values, key)
}

func (e *Env) Load() {
	newValues := make(map[string]string)
	err := database.TxQuery("SELECT key, value FROM environment WHERE env=?", e.name)(func(row *sql.Rows) error {
		var key string
		var val string
		if err := row.Scan(&key, &val); err != nil {
			return err
		}
		newValues[key] = val
		return nil
	})
	errored.Check(err, "env load")
	e.values = newValues
}

func (e *Env) Save(propagate bool) {
	err := database.Tx(func(tx *sql.Tx) error {
		if _, err := tx.Exec("DELETE FROM environment WHERE env=?", e.name); err != nil {
			return err
		}
		valueStmts := make([]string, 0, len(e.values))
		args := make([]any, 0, 3*len(e.values))
		for k, v := range e.values {
			valueStmts = append(valueStmts, "(?, ?, ?)")
			args = append(args, e.name)
			args = append(args, k)
			args = append(args, v)
		}
		query := "INSERT INTO environment(env, key, value) VALUES" + strings.Join(valueStmts, ",")
		_, err := tx.Exec(query, args...)
		return err
	})
	errored.Check(err, "env save")
	if propagate && e.parent != nil {
		e.parent.Save(true)
	}
}

func createEnvTable() {
	_, err := database.TxExec(
		`CREATE TABLE IF NOT EXISTS environment (
			env TEXT, key TEXT, value TEXT
		)`,
	)
	errored.Check(err, "init db.environment")
	Global.Load()
}
