package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cssivision/reverseproxy"
	"github.com/daafonsecato/kataterm-reverseproxy/pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	// Enable CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://terminal.kataterm.com")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// Handle OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	sc := handlers.NewSessionController()

	r.HandleFunc("/create", sc.createMachineHandler).Methods("GET")
	r.HandleFunc("/terminate", sc.terminateMachineHandler).Methods("POST", "OPTIONS")
	r.Use(corsMiddleware)

	http.ListenAndServe(":9090", r)

	proxy := &reverseproxy.ReverseProxy{
		Director:      handlers.customDirector(),
		Transport:     http.DefaultTransport,
		FlushInterval: 0,
		ErrorLog:      log.New(os.Stderr, "proxy: ", log.LstdFlags),
	}

	http.Handle("/", proxy)
	log.Println("Reverse proxy server started on :7070")
	log.Fatal(http.ListenAndServe(":7070", nil))
}
