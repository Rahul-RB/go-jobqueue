package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Rahul-RB/go-jobqueue/constants"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Nats struct {
	NatsConn *nats.Conn
	Stream   *jetstream.Stream
	Consumer *jetstream.Consumer
}

func CreateStreamAndConsumer() *Nats {
	nc, err := nats.Connect(nats.DefaultURL, nats.Name("Worker"))
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}

	ctx := context.Background()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create jetstream object:", err)
	}

	// Create one stream
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        constants.StreamName,
		Description: "Job Output Stream",
		Subjects:    []string{constants.SubjectName + ".*"},
	})
	if err != nil {
		log.Fatal("Failed to create a stream on jetstream:", err)
	}

	// Create one consumer
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    constants.ConsumerName,
		Durable: constants.ConsumerName,
	})
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}

	return &Nats{
		NatsConn: nc,
		Stream:   &stream,
		Consumer: &consumer,
	}

}

func InjectNats(n *Nats) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("nats", n)
		c.Next()
	}
}

func (n *Nats) Publish(ctx context.Context, s string, msg string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	js, err := jetstream.New(n.NatsConn)
	if err != nil {
		return errors.New("failed to create jetstream object:" + err.Error())
	}

	_, err = js.Publish(ctx, s, []byte(msg))
	if err != nil {
		return errors.New("failed to publish on" + s + ":" + err.Error())
	}
	return nil
}
