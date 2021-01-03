package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", smartHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func smartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	startTime := time.Now()

	done := make(chan struct{})
	go logRunningProcess(done)

	select {
	case <- done:
		fmt.Fprintln(w, "Hey there!")
	case <- ctx.Done():
		log.Println(ctx.Err())
		log.Println("after: ", time.Now().Sub(startTime).Seconds())
		http.Error(w, ctx.Err().Error(), http.StatusInternalServerError)
	}
}

func logRunningProcess(done chan struct{}) {
	time.Sleep(4 * time.Second)
	done <- struct{}{}
}
