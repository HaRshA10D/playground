package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/", nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Request processing error: ", err)
	}
	fmt.Println("response: ", res.StatusCode)
	io.Copy(os.Stdout, res.Body)
}
