package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wingedrhino/golang-snippets/networking/https-server-redis/handlers"
	"github.com/wingedrhino/golang-snippets/networking/https-server-redis/internal/platform/stack"
	"github.com/wingedrhino/golang-snippets/networking/https-server-redis/internal/util"
)

var redisListCap = flag.Int64("redis-list-cap", 15, "Capacity of list in Redis. Default value: '15'")
var redisListKey = flag.String("redis-list-key", "http-server-redis-requests", "Default value: 'http-server-redis-requests'")
var port = flag.String("port", ":8443", "Port to listen at; Defaults to ':8443'")
var redisURL = flag.String("redis-url", "localhost:6379", "Default value: 'localhost:6379'")
var redisPassword = flag.String("redis-password", "", "Default value: ''")
var redisDB = flag.Int("redis-db", 0, "Default value: '0'")
var certPath = flag.String("cert", "", "Path to the certificate file to use")
var keyPath = flag.String("key", "", "Path to private key of certificate to use.")

func initFlags() {
	flag.Parse()
	if len(*certPath) == 0 {
		fmt.Println("Argument cert is mandetory!")
		os.Exit(1)
	}
	if len(*keyPath) == 0 {
		fmt.Println("Argumnt key is mandetory!")
		os.Exit(1)
	}
	if len(*port) == 0 {
		fmt.Println("Argument port is empty. Using default value ':8443'.")
	}
}

func main() {
	initFlags()

	st, err := stack.NewRedisStack(*redisURL, *redisPassword, *redisDB, *redisListKey, *redisListCap)
	util.CheckFatal(err, "Unable to connect to Redis!")

	handler := handlers.NewHandler(st)
	tlsConfig := util.GetTLSConfig()

	h := &http.Server{
		Addr:      *port,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		err := h.ListenAndServeTLS(*certPath, *keyPath)
		if err != nil {
			// NOTE: Previously we'd check for err and call os.Exit()
			// immediately. Turns out that's a bad idea if you're waiting on a
			// context shutdown because ListenAndServeTLS is a blocking method
			// and can return error even if it successfully launched but
			// shutdown via context.
			fmt.Printf("Error in http.ListenAndServeTLS: %v\n\n", err)
		}
	}()

	<-stop

	fmt.Printf("\n\nShutting down the server. Time: %s\n\n", time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = h.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Error in shutting down server: %v\n", err)
	}
	fmt.Printf("Server gracefully stopped. Time: %s\n\n", time.Now())
}
