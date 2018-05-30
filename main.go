package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	var backend string
	flag.StringVar(&backend, "backend", "http://localhost:8444", "The base url of the backend")

	var virtualHost string
	flag.StringVar(&virtualHost, "host", "", "The virtual host to use for requests. e.g. connect-beta.nordu.net")

	var address string
	flag.StringVar(&address, "address", "127.0.0.1:8888", "The address to listen on")

	flag.Parse()

	p := &proxy{backend, virtualHost}

	server := &http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(p.proxyRequest),
	}

	log.Println("Starting proxy on http://" + address)
	log.Fatal(server.ListenAndServe())
}

type proxy struct {
	Backend     string
	VirtualHost string
}

func (i *proxy) proxyRequest(w http.ResponseWriter, req *http.Request) {
	target, err := url.Parse(i.Backend)
	if err != nil {
		log.Println("ERROR", "Unable to parse backen url:", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	req.URL.Host = target.Host
	req.URL.Scheme = target.Scheme
	if i.VirtualHost != "" {
		req.Host = i.VirtualHost
	}
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	// Bad practice when dealing with untrusted clients
	defer resp.Body.Close()
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
