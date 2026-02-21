package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// List of backend servers
	serverList := []string{":8081", ":8082"}
	pool := &ServerPool{}

	for _, s := range serverList {
		u, _ := url.Parse(s)

		pool.AddBackend(NewBackend(u))
	}

	// Start Health Check in background
	go func() {
		for {
			time.Sleep(20 * time.Second)
			pool.HealthCheck()
		}
	}()

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(lbHandler(pool)),
	}

	log.Println("Sentinel-LB started at :8080")
	log.Fatal(server.ListenAndServe())
}
