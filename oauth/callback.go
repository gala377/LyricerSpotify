package oauth

import (
	"context"
	"log"
	"net/http"
)

// ServCallback starts a http server waiting for oauth service callback.
//
// Server listens on http://localhost:{servAddr}{callback}
// for the callback request with parameter "code" passed in GET.
// Then sends the code value by the returned channel, closes it and
// shuts down the server.
//
// Example
//  codeChan := oauth.ServCallback("/callback", ":9090")
//  code := <-codeChan
func ServCallback(callbackRoute string, servAddr string) <-chan string {
	mux, closeChan, codeChan := setUpMux(callbackRoute)
	server := setUpServer(servAddr, mux)

	log.Printf("Set up server:\n\tport: %s\n\tcallback: %s\t\n ", servAddr, callbackRoute)
	log.Println("Starting callback server")
	go serveCallbackServer(server, closeChan)

	return codeChan
}

func setUpMux(callbackRoute string) (*http.ServeMux, chan bool, <-chan string) {
	shutdown := make(chan bool)
	code := make(chan string)

	mux := http.NewServeMux()
	mux.HandleFunc(callbackRoute, handleCallback(shutdown, code))

	return mux, shutdown, code
}

func setUpServer(servAddr string, mux http.Handler) *http.Server {
	serv := http.Server{
		Addr:    servAddr,
		Handler: mux,
	}
	return &serv
}

func handleCallback(closeChan chan bool, codeChan chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			code, ok := r.URL.Query()["code"]
			if !ok {
				log.Println("Code not present in callback uri")
				http.Error(w, "No code in uri", http.StatusBadRequest)
				return
			}
			log.Println("Authorizarion code received")
			codeChan <- code[0]
			closeChan <- true
		} else {
			log.Println("Request wasn't get")
			http.Error(w, "Wrong request", http.StatusBadRequest)
		}
	}
}

func serveCallbackServer(server *http.Server, closeChan <-chan bool) {
	go func() {
		err := server.ListenAndServe()
		log.Printf("Server closed, reason: %s\n", err)
	}()
	log.Println("Waiting for server shutdown")
	<-closeChan
	log.Println("Closing server")
	server.Shutdown(context.Background())
	log.Println("Server closed")
}
