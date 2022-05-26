package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func RunHttpServer(port string) error {
	mux := MakeRouter()

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Bypass-Tunnel-Reminder"},
	})

	log.Println("Listening on", port)
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        c.Handler(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
