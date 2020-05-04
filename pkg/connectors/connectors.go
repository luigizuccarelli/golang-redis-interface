package connectors

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connections struct {
	Http  *http.Client
	Redis *redis.Client
}

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

// NewClientConnectors returns Connectors struct
func NewClientConnections(logger *simple.Logger) Clients {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
	})

	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	conns := &Connections{Redis: redisClient, Http: httpClient}
	logger.Debug(fmt.Sprintf("Connection details %v\n", conns))
	return conns
}
