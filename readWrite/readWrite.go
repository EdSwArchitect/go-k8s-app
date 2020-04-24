package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var reading bool

func init() {
	flag.BoolVar(&reading, "reading", false, "If true, reads the /data directory. If false, writes random stuff to it")

	flag.Parse()
}

func tickTock() {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		name := uuid.New()

		filePath := fmt.Sprintf("/data/%s", name.String())

		f, err := os.Create(filePath)

		if err != nil {
			log.Fatalf("Unable to write file: %s\n", filePath)
		}

		_, err = f.WriteString(fmt.Sprintf("%s: %s\n", name.String(), time.Now()))

		f.Close()

		fmt.Println("Goober")

	}

}

func healthAndStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func getListing(w http.ResponseWriter, r *http.Request) {
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
	if reading {
		http.HandleFunc("/", healthAndStatus)
		http.HandleFunc("/getListing", getListing)
		log.Fatal(http.ListenAndServe(":8090", nil))

	} else {
		http.HandleFunc("/", healthAndStatus)
		http.HandleFunc("/getFile", getListing)
		log.Fatal(http.ListenAndServe(":8091", nil))

		ch := make(chan bool)

		go tickTock()

		<-ch
	}

}
