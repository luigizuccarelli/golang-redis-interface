package connectors

import (
	"net/http"
	"time"
)

// Client Interface - allows for different implmentations for testing and real environments
type Clients interface {
	Get(string) ([]byte, error)
	Set(string, []byte, time.Duration) error
	Del(string) error
	Do(req *http.Request) (*http.Response, error)
	Close() error
}
