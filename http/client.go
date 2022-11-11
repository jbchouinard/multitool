package http

import (
	"net/http"
	"time"
)

func MakeClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
