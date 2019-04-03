package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	// http.HandleFunc("/", indexHandler)
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)
	middleware := NewMiddleware(router, "I'm middleware")
	//http.Handle("/", router)
	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), middleware))
}

// The type of our middleware consists of the original handler we want to wrap and a message
type Middleware struct {
	next    http.Handler
	message string
}

// Make a constructor for our middleware type since its fields are not exported (in lowercase)
func NewMiddleware(next http.Handler, message string) *Middleware {
	return &Middleware{next: next, message: message}
}

// Our middleware handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We can modify the request here; for simplicity, we will just log a message
	log.Printf("msg: %s, Method: %s, URI: %s\n", m.message, r.Method, r.RequestURI)
	m.next.ServeHTTP(w, r)
	// We can modify the response here
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, world!")
}
func hello(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "Hello, world!", params.ByName("name"))
}
