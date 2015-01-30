package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var port = flag.Int("port", 3000, "Port to listen on")

func createBuild(w http.ResponseWriter, r *http.Request) {
	// Create build from params.
	io.WriteString(w, "createBuild")
}

func main() {
	flag.Parse()

	http.HandleFunc("/createBuild", createBuild)
	http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
}
