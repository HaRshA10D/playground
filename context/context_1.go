package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	doneChan := make(chan int, 1)

	go cancelWithTimeOut(doneChan)
	<- doneChan
	fmt.Println("example 1 done")

	go cancelImmediately(doneChan)
	<- doneChan
	fmt.Println("example 2 done")
}

func cancelWithTimeOut(done chan int) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()
	printAfter(ctx, 3 * time.Second, "example: cancel with timeout")

	done <- 0
}

func cancelImmediately(done chan int) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	printAfter(ctx, 3 * time.Second, "example: cancel manually")

	done <- 0
}

func printAfter(ctx context.Context, after time.Duration, message string) {
	select {
	case <- time.After(after):
		fmt.Println(message)
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
