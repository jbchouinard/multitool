package env

import (
	"errors"
	"strings"
)

type KV struct {
	Key   string
	Value string
}

func ParseKV(s string) (*KV, error) {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return nil, errors.New("expected key-value like foo=bar")
	}
	return &KV{parts[0], parts[1]}, nil
}

func ParseKVs(ss []string) ([]*KV, error) {
	params := make([]*KV, 0, len(ss))
	for _, s := range ss {
		if p, err := ParseKV(s); err != nil {
			return nil, err
		} else {
			params = append(params, p)
		}
	}
	return params, nil
}

func AddKVs(m map[string]string, kvs []*KV) {
	for _, p := range kvs {
		m[p.Key] = p.Value
	}
}
