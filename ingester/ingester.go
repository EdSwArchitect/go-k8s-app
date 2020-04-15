package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// AppConfig application configuration
type AppConfig struct {
	Port      int    `json:"port"`
	HealthURI string `json:"healthUri"`
	DirURI    string `json:"dirUri"`
	FilesDir  string `json:"filesDir"`
}

var configPath string
var theConfiguration AppConfig

func healthAndStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

}

func handleFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain")
	fmt.Fprintf(w, "Hello file")
}

func server() {
	http.HandleFunc("/", healthAndStatus)
	http.HandleFunc("/handleFile", handleFile)
}

func init() {
	log.Println("Ingester package init")

	flag.StringVar(&configPath, "configPath", "", "The configuration path")

	flag.Parse()

	log.Printf("The config path: '%s'\n", configPath)

	configData, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Printf("Config file error: %s.\n%s\n", configPath, err)
	}

	json.Unmarshal(configData, &theConfiguration)

	log.Printf("Configuration data: %+v\n", theConfiguration)

}

func main() {
	log.Println("Ingester main init")
	server()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", theConfiguration.Port), nil))

}
