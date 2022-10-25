package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var jobsFor3 = make(chan int, 50)
var jobsFor1 = make(chan int, 50)

func getGet1(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get1" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var i int
	json.NewDecoder(r.Body).Decode(&i)

	fmt.Printf("Received %d from First Server "+time.Now().Format(time.StampMilli)+"\n\n", i)
	jobsFor3 <- i
}

func getGet3(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get3" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var i int
	json.NewDecoder(r.Body).Decode(&i)

	fmt.Printf("Received %d from Third Server "+time.Now().Format(time.StampMilli)+"\n\n", i)
	jobsFor1 <- i
}

func HandleDataTo1(nr int) {
	for i := range jobsFor1 {
		orderMarshalled, _ := json.Marshal(i)
		responseBody := bytes.NewBuffer(orderMarshalled)

		http.Post("http://firstserver:8010/receive", "application/json", responseBody)

		//fmt.Printf("Thread %d sent data %d to First Server.\n\n", nr, i)
	}
}

func HandleDataTo3(nr int) {
	for i := range jobsFor3 {
		orderMarshalled, _ := json.Marshal(i)
		responseBody := bytes.NewBuffer(orderMarshalled)

		http.Post("http://thirdserver:8030/get", "application/json", responseBody)

		//fmt.Printf("Thread %d sent data to Third Server.\n\n", nr)
	}
}

func main() {
	http.HandleFunc("/get1", getGet1)
	http.HandleFunc("/get3", getGet3)

	for i := 1; i <= 3; i++ {
		go HandleDataTo3(i)
	}

	for i := 1; i <= 3; i++ {
		go HandleDataTo1(i)
	}

	fmt.Printf("Second Server started on PORT 8020\n")
	if err := http.ListenAndServe(":8020", nil); err != nil {
		log.Fatal(err)
	}

}
