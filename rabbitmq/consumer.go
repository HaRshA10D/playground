package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleErrorF(err)
	defer conn.Close()

	ch, err := conn.Channel()
	handleErrorF(err)

	declareTestQueueF(ch)

	msg, err := ch.Consume("test-queue", "", false, false, false, false, nil)
	handleErrorF(err)
	go startConsumer(ctx, msg)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	cancel()
}

func declareTestQueueF(ch *amqp.Channel) {
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
	handleErrorF(err)
	log.Print(testQueue)
}

func handleErrorF(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func startConsumer(ctx context.Context, msgChannel <-chan amqp.Delivery) {

	go func() {
		for msg := range msgChannel {
			log.Println("Message received")
			log.Println(string(msg.Body))
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Stopping consumption")
		return
	}
}
