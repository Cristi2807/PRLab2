package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var ID uint64

func incID() uint64 {
	return atomic.AddUint64(&ID, 1)
}

func getReceive(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/receive" {
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
}

func ProduceData(i int) {
	for {
		orderMarshalled, _ := json.Marshal(incID())
		responseBody := bytes.NewBuffer(orderMarshalled)

		http.Post("http://secondserver:8020/get1", "application/json", responseBody)

		//fmt.Printf("Thread %d sent data to Second Server.\n\n", i)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	http.HandleFunc("/receive", getReceive)

	for i := 1; i <= 6; i++ {
		go ProduceData(i)
	}

	fmt.Printf("First Server started on PORT 8010\n")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatal(err)
	}

}
