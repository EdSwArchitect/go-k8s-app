package main

import (
	"fmt"
	"log"
	"net/http"
)

func catchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain")
	fmt.Fprintln(w, "Default stuff")
}

func main() {
	http.HandleFunc("/", catchAll)

	log.Fatal(http.ListenAndServe(":1180", nil))
}
