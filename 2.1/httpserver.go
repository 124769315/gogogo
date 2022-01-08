package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", myHandler)
	mux.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":80", mux)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "********** go version **********\n")
	io.WriteString(w, runtime.Version()+"\n")
	io.WriteString(w, "********** request header **********\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))

	}

}
