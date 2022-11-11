package http

import "github.com/jbchouinard/wmt/env"

type APIDefinition struct {
	id      int64
	Name    string
	Env     *env.Env
	Headers map[string]string
}
