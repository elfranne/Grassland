package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	lookupJson := make(map[string]JsonItem)
	for _, i := range jsonContent {
		lookupJson[i.ID] = i
	}
	ipHost := strings.Split(lookupJson["final_score"].IP, "/")
	final_score := fmt.Sprintf("final_score{ip=\"%s\",url=\"%s\"} %s\n", ipHost[0], ipHost[1], lookupJson["final_score"].Finding)
	io.WriteString(w, final_score)
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
	fmt.Printf("Process started with PID: %d\n", cmd.Process.Pid)

}

func main() {
	fmt.Println("starting webservice")
	http.HandleFunc("/metrics", getMetrics)
	http.HandleFunc("/", processMetrics)
	err := http.ListenAndServe(":9232", nil)
	if err != nil {
		log.Fatal(err)
	}
}
