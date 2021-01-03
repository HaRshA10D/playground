package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "time"
)

func main() {
    fmt.Println("Hello world")

    waitChannel := make(chan os.Signal)
    printChannel := make(chan string, 1)

    go runLoop(printChannel)
    go printStream(printChannel)



    signal.Notify(waitChannel, os.Kill)
    <- waitChannel
}

func runLoop(printChannel chan string) {
    loop := 1
    for true {
        time.Sleep(2 * time.Second)
        printChannel <- fmt.Sprintf("Print loop: %d\n", loop)
        loop++

        if loop == 13 {
            close(printChannel)
            return
        }
    }
}

func printStream(printChannel chan string) {
    for val := range printChannel {
        fmt.Printf(val)
    }
    fmt.Printf("Closing print channel\n")
}

func runGreeting(printChannel chan string, sequence int) {

    url := fmt.Sprintf("http://localhost/greetings?sequence=%d", sequence)
    http.Get(url)
}
