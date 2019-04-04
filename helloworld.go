package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"time"
)

const username = "gordon"
const pass = "secret!"

type Post struct {
	Title  string
	Teaser string
	Posted time.Time
	URL    string
	Tag    string
}

func main() {

	// http.HandleFunc("/", indexHandler)
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)
	router.POST("/posts", postPost)
	router.GET("/protected/", BasicAuth(Protected))
	http.Handle("/", router)
	appengine.Main()

}
func Protected(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Protected!\n")
}

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == username && password == pass {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello, world!")
}

func postPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx := appengine.NewContext(r)
	post := Post{
		Title:  r.FormValue("title"),
		Teaser: r.FormValue("teaser"),
		Tag:    r.FormValue("tag"),
		URL:    r.FormValue("url"),
		Posted: time.Now(),
	}
	key := datastore.NewIncompleteKey(ctx, "Post", nil)

	if completeKey, err := datastore.Put(ctx, key, &post); err != nil {
		log.Errorf(ctx, "datastore.Put: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	} else {
		log.Errorf(ctx, "delay Call %v", completeKey)
	}

}
func hello(w http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	fmt.Fprint(w, "Hello, world!", params.ByName("name"))

}
