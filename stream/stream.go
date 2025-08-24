package stream

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Rahul-RB/go-jobqueue/constants"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Stream struct {
	natsConn      *nats.Conn
	natsJetStream *jetstream.JetStream
	natsStream    jetstream.Stream
}

type StreamSession struct {
	isClosed     bool
	jobId        string
	natsConsumer jetstream.Consumer
	wsConn       *websocket.Conn
	messages     chan []byte
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

	return &Stream{
		natsConn:      nc,
		natsJetStream: &js,
		natsStream:    stream,
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

	if _, err := (*s.natsJetStream).Publish(ctx, sub, []byte(msg)); err != nil {
		return errors.New("failed to publish on" + sub + ":" + err.Error())
	}
	return nil
}

func (s *StreamSession) Read() {
	var cctx jetstream.ConsumeContext
	cctx, err := s.natsConsumer.Consume(func(message jetstream.Msg) {
		data := message.Data()
		if err := message.Ack(); err != nil {
			log.Println(s.jobId, "failed to ack:", err)
			return
		}
		if s.isClosed {
			log.Println(s.jobId, "has closed")
			s.isClosed = true
			close(s.messages)
			cctx.Stop()
			s.wsConn.Close()
			return
		}
		s.messages <- data
	})
	if err != nil {
		log.Println(s.jobId, "failed to consume:", err)
		return
	}
}

func (s *StreamSession) Write() {
	defer func() {
		log.Println(s.jobId, "job consumer closed via write")
		s.isClosed = true
		s.wsConn.Close()
	}()

	for message := range s.messages {
		log.Println(s.jobId, "got message:", string(message))
		if s.isClosed {
			return
		}
		if s.wsConn.WriteMessage(websocket.TextMessage, message) != nil {
			return
		}
	}
}

func (s *Stream) StartNewConsumer(id string, w *websocket.Conn) error {
	ctx := context.Background()
	consumerName := fmt.Sprintf("%v_%v", "job_output_stream_consumer", time.Now().Unix())
	consumer, err := s.natsStream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    consumerName,
		Durable: consumerName,
	})
	if err != nil {
		return err
	}

	session := &StreamSession{
		isClosed:     false,
		jobId:        id,
		natsConsumer: consumer,
		wsConn:       w,
		messages:     make(chan []byte, 256),
	}

	go session.Read()
	go session.Write()

	return nil
}
