// +build real

package connectors

import (
	"fmt"
	"net/http"
	"time"
)

// redis Get
func (r *Connections) Get(key string) ([]byte, error) {
	val, ok := r.Redis.Get(key).Result()
	if ok != nil {
		return nil, fmt.Errorf("Error: %s", ok)
	}
	return []byte(val), nil
}

// redis Set
func (r *Connections) Set(key string, value []byte, expr time.Duration) error {
	err := r.Redis.Set(key, value, 0).Err()
	return err
}

// redis Del
func (r *Connections) Del(key string) error {
	err := r.Redis.Del(key).Err()
	return err
}

// redis Close
func (r *Connections) Close() error {
	return nil
}

func (r *Connections) Do(req *http.Request) (*http.Response, error) {
	return r.Http.Do(req)
}
