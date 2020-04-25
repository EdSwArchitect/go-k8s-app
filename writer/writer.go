package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var reading bool

func init() {
	fmt.Println("Writer starting up.")
}

func tickTock() {
	ticker := time.NewTicker(15 * time.Second)

	for range ticker.C {
		name := uuid.New()

		filePath := fmt.Sprintf("/data/%s", name.String())

		f, err := os.Create(filePath)

		log.Printf("Writing file: %s\n", filePath)

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

func main() {
	http.HandleFunc("/", healthAndStatus)
	log.Fatal(http.ListenAndServe(":8091", nil))

	ch := make(chan bool)

	go tickTock()

	<-ch

}
