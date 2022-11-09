package env

import (
	"errors"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

func ParseParam(s string) (*Param, error) {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return nil, errors.New("expected param like foo=bar")
	}
	return &Param{parts[0], parts[1]}, nil
}

func ParseParams(ss []string) ([]*Param, error) {
	params := make([]*Param, 0, len(ss))
	for _, s := range ss {
		if p, err := ParseParam(s); err != nil {
			return nil, err
		} else {
			params = append(params, p)
		}
	}
	return params, nil
}

func AddParams(m map[string]string, params []*Param) {
	for _, p := range params {
		m[p.Key] = p.Value
	}
}
