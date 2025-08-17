package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/Rahul-RB/go-jobqueue/constants"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Manager"))
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	ctx := context.Background()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create jetstream object:", err)
	}

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        constants.StreamName,
		Description: "Job stream",
		Subjects:    []string{constants.SubjectName},
	})
	if err != nil {
		log.Fatal("Failed to create a stream on jetstream:", err)
	}

	_, err = stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    constants.ConsumerName,
		Durable: constants.ConsumerName,
	})
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}

	_, err = js.Publish(ctx, constants.SubjectName, []byte("Hello World from ts1!"))
	if err != nil {
		log.Println("Failed to publish on job_stream:", err)
	}

}
