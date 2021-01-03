package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err)

	declareTestQueue(ch)

	go startPublisher(ctx, ch)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()
}

func declareTestQueue(ch *amqp.Channel) {
	/*
		-------------
		Every queue declared gets a default binding to the empty exchange "" which has
		the type "direct" with the routing key matching the queue's name.

		-------------
		Durable and Non-Auto-Deleted -> Survive server restart and lives when empty too
		Non-Durable and Non-Auto-Deleted -> Lives when empty but will be deleted when server restart
		Durable and Auto-Deleted -> Survives broker restart but will bve deleted on empty

		-------------

	*/
	testQueue, err := ch.QueueDeclare("test-queue", true, false, false, false, nil)
	handleError(err)
	log.Print(testQueue)
}

func handleError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func startPublisher(context context.Context, ch *amqp.Channel) {

	count := 1
	for {
		select {
		case <-context.Done():
			log.Println("Done with periodic publishing")
			return
		case <-time.After(3 * time.Second):
			message := amqp.Publishing{
				ContentType:  "json/application",
				DeliveryMode: amqp.Persistent,
				Body:         []byte("{\"count\": \"" + strconv.Itoa(count) + "\"}"),
			}
			publishMessage(ch, message)
			count++
			log.Println("Published message no: " + strconv.Itoa(count))
		}
	}
}

func publishMessage(ch *amqp.Channel, msg amqp.Publishing) {
	/*

		Publishing is asynchronous, so listen to NotifyReturn of the channel for undelivered messages
	*/
	err := ch.Publish("", "test-queue", true, false, msg)
	handleError(err)
}
