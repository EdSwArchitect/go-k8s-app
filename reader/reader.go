package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var reading bool

func init() {
	fmt.Println("Reading starting up")
}

func healthAndStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func getListing(w http.ResponseWriter, r *http.Request) {

	fmt.Println("getListing")

	files, err := ioutil.ReadDir("/data")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unable to read directory /data")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain")

	for i, file := range files {
		fmt.Fprintf(w, "%d: %s\n", i, file.Name())
	}

}

func getFile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("getFile")

	url := r.URL

	requestURI := url.RequestURI()

	fmt.Fprintf(w, "Request URI: %s\n", requestURI)

	values := url.Query()

	fmt.Fprintf(w, "The values: %+v\n", values)

	filePath, ok := values["file"]

	if !ok || len(filePath[0]) == 0 {
		log.Println("Query parameter error")
		return
	}

	fmt.Fprintf(w, "The filePath: %s\n", filePath[0])

	theFile := fmt.Sprintf("/data/%s", filePath[0])

	file, err := os.Open(theFile)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("content-type", "text/plain")
		fmt.Fprintf(w, "Unable to get file: %s. %+v\n", filePath[0], err)

	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", scanner.Text())
	}

	file.Close()

}

func main() {
	http.HandleFunc("/", healthAndStatus)
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/listing", getListing)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
