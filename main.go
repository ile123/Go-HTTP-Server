package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"ile123/http-server/constants"
	"ile123/http-server/handlers"
	"ile123/http-server/middlewares"
	"net"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.GetRoot)
	mux.HandleFunc("/hello", handlers.GetHello)
	mux.HandleFunc("/api/message", handlers.HandleJSON)
	mux.HandleFunc("/health", handlers.HealthCheck)

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	stack := middlewares.NewStack()
	stack.Use(middlewares.LoggingMiddleware)
	stack.Use(middlewares.HeaderMiddleware)
	stack.Use(middlewares.SecurityHeadersMiddleware)

	handler := stack.Then(mux)

	ctx, cancelCtx := context.WithCancel(context.Background())

	newServer := func(addr string) *http.Server {
		return &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSConfig:    tlsConfig,
			BaseContext: func(l net.Listener) context.Context {
				return context.WithValue(ctx, constants.KeyServerAddr, l.Addr().String())
			},
		}
	}

	run := func(s *http.Server, name string) {
		//err := s.ListenAndServeTLS("cert.pem", "key.pem")
		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("%s closed\n", name)
		} else if err != nil {
			fmt.Printf("error listening for %s: %s\n", name, err)
		}
		cancelCtx()
	}

	go run(newServer(":3333"), "server one")
	go run(newServer(":4444"), "server two")

	<-ctx.Done()
}
