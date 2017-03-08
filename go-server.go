package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//all requests to top level pages are stored in ./content folder
func landHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if path == "" || path == "index" {
		path = "index.html"
	}
	//open file and send
	f, err := os.Open("content/" + path)
	if err != nil {
		fourOhFour(w, r)
	} else {
		http.ServeContent(w, r, path, time.Now(), f)
	}
}

// /var/www/ear7h-net/path folder contains user content
func pathHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/www/ear7h-net/" + r.URL.Path[1:]
	//open file and send
	f, err := os.Open(path)
	if err != nil {
		fourOhFour(w, r)
	} else {
		http.ServeContent(w, r, r.URL.Path, time.Now(), f)
	}
}

// /var/www/ear7h-net/bin folder contains heavy media, videos, sound, etc.
func binHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/www/ear7h-net/" + r.URL.Path[1:]
	//open file and send
	f, err := os.Open(path)
	if err != nil {
		fourOhFour(w, r)
	} else {
		http.ServeContent(w, r, r.URL.Path, time.Now(), f)
	}
}

//proxies requests to nodejs api server
/*
func apiHandler(w http.ResponseWriter, r *http.Request) {
	host := &url.URL{
		Scheme: "http",
		Host:   "localhost:81",
	}
	httputil.NewSingleHostReverseProxy(host).ServeHTTP(w, r)
}
*/

// proxies requests to localhost on port specified after /localproxy/
//ie. host.com/localproxy/8080/index.html goes to localhost:8080/index.html

func fwdlocalHandler(w http.ResponseWriter, r *http.Request) {
	rpath := strings.Split(string(r.URL.Path), "/")
	//only forward to ports 8000 - 8100 else 404
	it, err := strconv.Atoi(rpath[2])
	if err != nil {
		fourOhFour(w, r)
	}
	if it < 8000 || it > 8100 {
		fourOhFour(w, r)
		return
	}
	//concatination of url
	fwdURL := "http://localhost:" + rpath[2] + "/" + strings.Join(rpath[3:], "/")
	//send get request
	resp, err := http.Get(fwdURL)
	if err != nil {
		errLog(err)
		return
	}
	defer resp.Body.Close()
	//put body contents into body variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errLog(err)
		return
	}
	//send request
	fmt.Fprint(w, string(body))
}

func fourOhFour(w http.ResponseWriter, r *http.Request) {
	four04, _ := os.Open("content/notfound.html")
	http.ServeContent(w, r, "", time.Now(), four04)
}

func errLog(e error) {
	t := time.Now().String()
	f, err := os.OpenFile("err.log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString("$ " + t + "\n" + err.Error() + "\n\n"); err != nil {
		panic(err)
	}
}

func main() {
	finish := make(chan bool)
	fmt.Print("\nSERVER STARTING\n\n")

	server80 := http.NewServeMux()

	server80.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://ear7h.net/", 308)
	})

	server443 := http.NewServeMux()

	server443.HandleFunc("/", landHandler)
	server443.HandleFunc("/users/", pathHandler)
	server443.HandleFunc("/bin/", binHandler)
	server443.HandleFunc("/fwdlocal/", fwdlocalHandler)
	//erver443.HandleFunc("/api/", apiHandler)

	go func() {
		fmt.Println("server running on :80")
		e := http.ListenAndServe(":80", server80)
		fmt.Println(e)
	}()
	go func() {
		fmt.Println("server running on :443")
		/*
			e := http.ListenAndServeTLS(":443",
				"/etc/letsencrypt/live/ear7h.net/cert.pem",
				"/etc/letsencrypt/live/ear7h.net/privkey.pem",
				server443)
		*/e := http.ListenAndServe(":443", server443)
		fmt.Println(e)
	}()

	<-finish
}
