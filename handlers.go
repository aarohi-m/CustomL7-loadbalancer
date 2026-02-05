package main

import "net/http"

func lbHandler(pool *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := pool.GetNext()
		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}
