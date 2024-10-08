package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type JsonItem struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Severity string `json:"severity"`
	Finding  string `json:"finding"`
}

var jsonContent []JsonItem

func update() {
	// clear old data
	jsonContent = nil

	jsonRead, err := ioutil.ReadFile("./testssl.json")
	if err != nil {
		log.Println("Error when opening file: ", err)
	}

	err = json.Unmarshal(jsonRead, &jsonContent)
	if err != nil {
		log.Println("Error during Unmarshal(): ", err)
	}
	for _, i := range jsonContent {
		fmt.Println(i.ID)
	}

}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	for _, i := range jsonContent {
		io.WriteString(w, i.ID)
	}
}

func main() {

	http.HandleFunc("/metrics", getMetrics)
	err := http.ListenAndServe(":9232", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if strings.HasSuffix(event.Name, "cron.lock") {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					fmt.Printf("%s %s\n", event.Op, event.Name)
					update()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
}
