package main

import "net/http"

func main() {

	servMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: servMux,
	}
	server.ListenAndServe()
}
