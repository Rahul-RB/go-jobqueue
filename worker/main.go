package main

import (
	"context"
	"log"
	"sync"

	"github.com/Rahul-RB/go-jobqueue/constants"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Worker"))
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	ctx := context.Background()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create jetstream object:", err)
	}

	stream, err := js.Stream(ctx, constants.StreamName)
	if err != nil {
		log.Fatal("Failed to connect to test_stream:", err)
	}

	consumer, err := stream.Consumer(ctx, constants.ConsumerName)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}

	var wg sync.WaitGroup
	cctx, err := consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Received subject: %v message: %v", msg.Subject(), string(msg.Data()))
		if err := msg.Ack(); err != nil {
			log.Fatal("Failed to ack:", err)
		}
	})
	if err != nil {
		log.Fatal("Failed to consume:", err)
	}
	defer cctx.Stop()
	wg.Add(1)
	wg.Wait()
}
