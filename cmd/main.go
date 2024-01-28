package main

import (
	"github.com/david8128/kataterm-reverseproxy/internal/app/handlers"
	"github.com/david8128/kataterm-reverseproxy/internal/database"
	"github.com/david8128/kataterm-reverseproxy/pkg/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	sc := controllers.NewSessionController()

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		createMachineHandler(db, w, r)
	}).Methods("GET")

	http.HandleFunc("/terminate", func(w http.ResponseWriter, r *http.Request) {
		terminateMachineHandler(db, w, r)
	})
	r.Use(corsMiddleware)

	http.ListenAndServe(":9090", r)

	proxy := &reverseproxy.ReverseProxy{
		Director:      customDirector(db),
		Transport:     http.DefaultTransport,
		FlushInterval: 0,
		ErrorLog:      log.New(os.Stderr, "proxy: ", log.LstdFlags),
	}

	http.Handle("/", proxy)
	log.Println("Reverse proxy server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
