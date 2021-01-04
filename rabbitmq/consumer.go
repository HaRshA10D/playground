package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	/*
	----------------
	All deliveries in AMQP must be acknowledged.  It is expected of the consumer to
	call Delivery.Ack after it has successfully processed the delivery.  If the
	consumer is cancelled or the channel or connection is closed any unacknowledged
	deliveries will be requeued at the end of the same queue.
	----------------

	 */
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
			log.Println("Message received:" + string(msg.Body))
			go process(msg)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Stopping consumption")
		return
	}
}

func process(msg amqp.Delivery) {
	<-time.After(3 * time.Second)
	log.Println("Message Processed:" + string(msg.Body))
	msg.Ack(false)
}
