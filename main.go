package main

import "net/http"

//type apiHandler struct{}
//
//func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
//
//}

func main() {

	servMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: servMux,
	}
	servMux.Handle("/assets", http.FileServer(http.Dir("./assets")))
	servMux.Handle("/", http.FileServer(http.Dir(".")))

	server.ListenAndServe()
}
