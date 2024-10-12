package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type JsonItem struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Severity string `json:"severity"`
	Finding  string `json:"finding"`
}

var jsonContent []JsonItem
var idArray []string

func getMetrics(w http.ResponseWriter, r *http.Request) {
	jsonContent = nil
	jsonRead, err := ioutil.ReadFile("./out.json")
	if err != nil {
		log.Println("Error when opening file: ", err)
	}

	err = json.Unmarshal(jsonRead, &jsonContent)
	if err != nil {
		log.Println("Error during Unmarshal(): ", err)
	}
	jsonLookup := make(map[string]JsonItem)
	for _, i := range jsonContent {
		jsonLookup[i.ID] = i
	}
	for _, i := range idArray {
		ipHost := strings.Split(jsonLookup[string(i)].IP, "/")
		fmt.Fprintf(w, "final_score{ip=\"%s\",url=\"%s\"} %s\n", ipHost[1], ipHost[0], jsonLookup[string(i)].Finding)
	}
}

func processMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the process in the background
	cmd := exec.CommandContext(ctx, "./cron.sh")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	log.Printf("Process started with PID: %d\n", cmd.Process.Pid)
}

func main() {
	//  set CLI args
	address := flag.String("address", ":9232", "port or address to listen to, default to: :9232")
	id := flag.String("id", "final_score", "list of space seperated ID to expose, default to: final_score")
	flag.Parse()
	idArray = strings.Split(*id, " ")
	// start webserver
	log.Printf("starting webservice on %s, parsing %s", *address, strings.Split(*id, " "))
	http.HandleFunc("/metrics", getMetrics)
	http.HandleFunc("/", processMetrics)
	err := http.ListenAndServe(*address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
