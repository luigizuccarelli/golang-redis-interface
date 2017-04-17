package main 
import (
    "net/http"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func startHttpServer(port string) *http.Server {
    srv := &http.Server{Addr: ":" + port}

    http.HandleFunc("/", SimpleHandler)
    http.HandleFunc("/isalive", IsAlive)

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }
    }()

    return srv
}

func main() {
    srv := startHttpServer("9000")
    log.Printf("main: starting server on port %s", srv.Addr)
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
    log.Printf("main: server shutdown successfully")
    os.Exit(code)
}
