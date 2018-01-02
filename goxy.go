package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"net/url"
)

var (
	port = flag.Int("port", 8080, "Port number to start proxy on")
)

func init() {
	flag.Parse()
}

// https://github.com/Gonzih/http-forward-proxy/main.go
func copyHeaders(dst *http.Header, src *http.Header) {
	for k, vals := range *src {
		for _, v := range vals {
			if k == "goxy-url" {
				continue
			}
			log.Printf("Copying header %s: %s", k, v)
			dst.Set(k, v)
		}
	}
}

func newDirector(r *http.Request) func(*http.Request) {
	return func(req *http.Request) {

		u, err := url.Parse(r.Header.Get("goxy-url"))
		if err != nil {
			log.Println(err)
		}
		req.URL = u

		log.Println("New:", u.String())
		log.Println("New:", req.URL.String())


		reqLog, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			log.Printf("Got error %s\n %+v\n", err.Error(), req)
		}

		log.Println(string(reqLog))
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	proxy := &httputil.ReverseProxy{
		Transport: &http.Transport{},
		Director:  newDirector(r),
	}
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
