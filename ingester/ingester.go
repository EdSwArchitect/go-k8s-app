package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
)

// AppConfig application configuration
type AppConfig struct {
	Port      int    `json:"port"`
	HealthURI string `json:"healthUri"`
	DirURI    string `json:"dirUri"`
	FilesDir  string `json:"filesDir"`
	Elastic   string `json:"elastic"`
	Index     string `json:"index"`
}

var theConfiguration AppConfig
var es *elasticsearch.Client
var configPath string

func elastic(w http.ResponseWriter, r *http.Request) {

	// cfg := elasticsearch.Config{
	// 	Addresses: []string{
	// 		theConfiguration.Elastic,
	// 	},
	// }

	// es, err := elasticsearch.NewClient(cfg)

	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Header().Add("content-type", "text/plain")
	// 	fmt.Fprintf(w, "Unable to connect to Elastic. %+v\n", err)
	// 	return
	// }

	// res, err := es.Info()

	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	w.Header().Add("content-type", "text/plain")
	// 	fmt.Fprintf(w, "Unable to get Elastic info: %+v\n", err)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	// w.Header().Add("content-type", "text/plain")
	// fmt.Fprintf(w, "%+v", res)

}

func healthAndStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/plain")

	url := r.URL

	requestURI := url.RequestURI()

	fmt.Fprintf(w, "Request URI: %s\n", requestURI)

	values := url.Query()

	fmt.Fprintf(w, "The values: %+v\n", values)

	filePath, ok := values["filePath"]

	if !ok || len(filePath[0]) == 0 {
		log.Println("Query parameter error")
		return
	}

	fmt.Fprintf(w, "The filePath: %s\n", filePath[0])

	theFile := fmt.Sprintf("%s/%s", theConfiguration.FilesDir, filePath[0])

	// file, err := os.Open(filePath[0])

	file, err := os.Open(theFile)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Add("content-type", "text/plain")
		fmt.Fprintf(w, "Unable to get file: %s. %+v\n", filePath[0], err)

	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	load(txtlines)

	// fmt.Fprintf(w, "Helloload file")
}

// load the data into ElasticSearch
func load(lines []string) {
	var wg sync.WaitGroup

	for _, line := range lines {
		wg.Add(1)

		go func(line string) {
			defer wg.Done()

			var b strings.Builder
			b.WriteString(`{"line" : "`)
			b.WriteString(line)
			b.WriteString(`"}`)

			idx, _ := uuid.NewRandom()

			log.Printf("UUID: %s\n", idx.String())
			log.Printf("%s\n", b.String())

			req := esapi.IndexRequest{
				Index:      theConfiguration.Index,
				DocumentID: idx.String(),
				Body:       strings.NewReader(b.String()),
				Refresh:    "true",
			}

			if es == nil {
				log.Fatal("es is now null!")
			}

			res, err := req.Do(context.Background(), es)

			if err != nil {
				log.Printf("Error")
				return
			}

			defer res.Body.Close()

			if res.IsError() {
				return
			} else {
				var r map[string]interface{}

				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
					return
				} else {
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}

		}(line)
	}

	wg.Wait()

}

func server() {
	http.HandleFunc("/", healthAndStatus)
	http.HandleFunc("/handleFile", handleFile)
	// http.HandleFunc("/elastic", elastic)

}

func init() {
	log.Println("Ingester package init")

	flag.StringVar(&configPath, "configPath", "", "The configuration path")

	flag.Parse()

	log.Printf("The config path: '%s'\n", configPath)

	if len(configPath) == 0 {
		configPath = "/data/config.json"
	}

	configData, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Printf("Config file error: %s.\n%s\n", configPath, err)
	}

	json.Unmarshal(configData, &theConfiguration)

	log.Printf("Configurations: %+v\n", theConfiguration)

	// connectES()

}

func connectES() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			theConfiguration.Elastic,
		},
	}

	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Unable to connect to Elastic. %+v\n", err)
	}

	res, err := es.Info()

	if err != nil {
		log.Fatalf("Unable to get Elastic info: %+v\n", err)
	}

	log.Printf("Configuration data: %+v\n", theConfiguration)
	log.Printf("Elastic\n%+v\n", res)

}

func main() {
	log.Println("Ingester main init")
	server()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", theConfiguration.Port), nil))

}
