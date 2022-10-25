package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var jobs = make(chan int, 50)

func getGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var i int
	json.NewDecoder(r.Body).Decode(&i)

	fmt.Printf("Received %d from Second Server "+time.Now().Format(time.StampMilli)+"\n\n", i)
	jobs <- i
}

func HandleJobs(nr int) {
	for i := range jobs {
		orderMarshalled, _ := json.Marshal(i)
		responseBody := bytes.NewBuffer(orderMarshalled)

		time.Sleep(1 * time.Second)

		http.Post("http://secondserver:8020/get3", "application/json", responseBody)

		//fmt.Printf("Thread %d sent data %d to Third Server.\n\n", nr, i)
	}
}

func main() {
	http.HandleFunc("/get", getGet)

	for i := 1; i <= 3; i++ {
		go HandleJobs(i)
	}

	fmt.Printf("Third Server started on PORT 8030\n")
	if err := http.ListenAndServe(":8030", nil); err != nil {
		log.Fatal(err)
	}

}
