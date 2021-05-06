// +build real

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
