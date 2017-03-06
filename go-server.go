package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

//all requests to top level pages are stored in ./content folder
func landHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "" || path == "index" {
		path = "index.html"
	}

	f, err := os.Open("content/" + path)
	if err != nil {
		four04, _ := os.Open("content/notfound.html")
		http.ServeContent(w, r, "", time.Now(), four04)
	} else {
		http.ServeContent(w, r, path, time.Now(), f)
	}
}

// /var/www/ear7h-net/path folder contains user content
func pathHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/www/ear7h-net/" + r.URL.Path[1:]
	f, err := os.Open(path)
	if err != nil {
		four04, _ := os.Open("content/notfound.html")
		http.ServeContent(w, r, "", time.Now(), four04)
	} else {
		http.ServeContent(w, r, r.URL.Path, time.Now(), f)
	}
}

// /var/www/ear7h-net/bin folder contains heavy media, videos, sound, etc.
func binHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/www/ear7h-net/" + r.URL.Path[1:]
	f, err := os.Open(path)
	if err != nil {
		four04, _ := os.Open("content/notfound.html")
		http.ServeContent(w, r, "", time.Now(), four04)
	} else {
		http.ServeContent(w, r, r.URL.Path, time.Now(), f)
	}
}

/*
func apiHandler(w http.ResponseWriter, r *http.Request) {
}
*/

// proxies requests to localhost on port specified after /localproxy/
//ie. host.com/localproxy/8080/index.html goes to localhost:8080/index.html

func localproxyHandler(w http.ResponseWriter, r *http.Request) {
	host := "http://localhost:" + r.URL.Path[1] + "/" + r.URL.Path[2:]
	httputil.NewSingleHostReverseProxy(host).ServeHTTP(w, r)
}

func main() {
	finish := make(chan bool)
	fmt.Println("\nSERVER STARTING\n")

	server80 := http.NewServeMux()

	server80.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://ear7h.net/", 308)
	})

	server443 := http.NewServeMux()

	server443.HandleFunc("/", landHandler)
	server443.HandleFunc("/users/", pathHandler)
	server443.HandleFunc("/bin/", binHandler)
	server443.Handle("/localproxy/", localproxyHandler)
	//server443.HandleFunc("/api/", apiHandler)

	go func() {
		fmt.Println("server running on :80")
		e := http.ListenAndServe(":80", server80)
		fmt.Println(e)
	}()
	go func() {
		fmt.Println("server running on :443")
		e := http.ListenAndServeTLS(":443",
			"/etc/letsencrypt/live/ear7h.net/cert.pem",
			"/etc/letsencrypt/live/ear7h.net/privkey.pem",
			server443)
		//e := http.ListenAndServe(":443", server443)
		fmt.Println(e)
	}()

	<-finish
}
