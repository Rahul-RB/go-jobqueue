package stream

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

type Stream struct {
	natsConn      *nats.Conn
	natsStreamObj *jetstream.JetStream
	natsConsumer  *jetstream.Consumer
}

func NewStream() *Stream {
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

	return &Stream{
		natsConn:      nc,
		natsStreamObj: &js,
		natsConsumer:  &consumer,
	}

}

func InjectStream(s *Stream) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("stream", s)
		c.Next()
	}
}

func (s *Stream) Publish(ctx context.Context, sub string, msg string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if _, err := (*s.natsStreamObj).Publish(ctx, sub, []byte(msg)); err != nil {
		return errors.New("failed to publish on" + sub + ":" + err.Error())
	}
	return nil
}
