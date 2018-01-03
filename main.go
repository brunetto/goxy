package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"net/url"
	"io"
	"strings"
)

var (
	port = flag.Int("port", 8080, "Port number to start proxy on")
)

func init() {
	flag.Parse()
}



//func newDirector(r *http.Request) func(*http.Request) {
//	return func(req *http.Request) {
//
//		u, err := url.Parse(r.Header.Get("goxy-url"))
//		if err != nil {
//			log.Println(err)
//		}
//		req.URL = u
//
//		log.Println("New:", u.String())
//		log.Println("New:", req.URL.String())
//
//
//		reqLog, err := httputil.DumpRequestOut(req, false)
//		if err != nil {
//			log.Printf("Got error %s\n %+v\n", err.Error(), req)
//		}
//
//		log.Println(string(reqLog))
//	}
//}
//
//func proxyHandler(w http.ResponseWriter, r *http.Request) {
//	proxy := &httputil.ReverseProxy{
//		Transport: &http.Transport{},
//		Director:  newDirector(r),
//	}
//	proxy.ServeHTTP(w, r)
//}

// https://github.com/Gonzih/http-forward-proxy/main.go
func copyHeaders(dst *http.Header, src *http.Header, direction string ) {
	for k, vals := range *src {
		for _, v := range vals {
			if strings.ToLower(k) == "goxy-url" {
				continue
			}
			log.Printf("Copying %v header %s: %s", direction, k, v)
			dst.Set(k, v)
		}
	}
}

func proxyHandler (w http.ResponseWriter, r *http.Request) {
	method := r.Method
	u, err := url.Parse(r.Header.Get("goxy-url"))
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest(method, u.String(), r.Body)
	if err != nil {
		log.Printf("Error while creating request %s\n", err)
		return
	}

	sourceHeaders := r.Header
	destinationHeaders := req.Header
	copyHeaders(&destinationHeaders, &sourceHeaders, "outgoing")
	proxyResp, err := client.Do(req)

	if err != nil {
		log.Printf("Error while executing request %s\n", err)
		return
	}

	sourceHeaders = proxyResp.Header
	destinationHeaders = w.Header()
	copyHeaders(&destinationHeaders, &sourceHeaders, "incoming")

	io.Copy(w, proxyResp.Body)

}

func main() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
