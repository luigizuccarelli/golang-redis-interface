package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger  Logger
	counter uint64
)

func startHttpServer(port string) *http.Server {
	logger.Level = "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logger.Level = os.Getenv("LOG_LEVEL")
	}
	srv := &http.Server{Addr: ":" + port}

	http.HandleFunc("/simulator", SimpleHandler)
	http.HandleFunc("/isalive", IsAlive)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {
	var port string = "9000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	srv := startHttpServer(port)
	logger.Info("Starting server on port " + srv.Addr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}
