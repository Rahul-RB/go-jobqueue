package jobs

import (
	"bufio"
	"context"
	"errors"
	"log"
	"sync"

	"github.com/Rahul-RB/go-jobqueue/constants"
	"github.com/Rahul-RB/go-jobqueue/stream"
	"github.com/Rahul-RB/go-jobqueue/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Job struct {
	Id     string         `json:"id"`
	Name   string         `json:"name"`
	stream *stream.Stream `json:"-"`
}

var jobs map[string]*Job = make(map[string]*Job)
var jobLock sync.Mutex

func GetJob(id string) (*Job, error) {
	jobLock.Lock()
	defer jobLock.Unlock()
	if j, ok := jobs[id]; ok {
		return j, nil
	}
	return &Job{}, errors.New("Couldn't find job with ID:" + id)
}

func (j *Job) Run() {
	// run job
	// get output of job
	// publish that to Stream
	ctx := context.Background()
	cmd, err := utils.RunWithTimeout("./dummy-job/dummy-job", "-name", j.Name, "-interval", "1s")
	if err != nil {
		log.Fatal("Failed to run command:", err.Error())
		if err := j.stream.Publish(ctx, j.Name, err.Error()); err != nil {
			log.Fatal("Failed to publish to:", j.Name, err.Error())
		}
	}

	scanner := bufio.NewScanner(*cmd.Stdout)
	for scanner.Scan() {
		if err := j.stream.Publish(ctx, j.Name, scanner.Text()); err != nil {
			log.Fatal("Failed to publish to:", j.Name, err.Error())
		}
	}

	cmd.Cmd.Wait()
}

func NewJob(s *stream.Stream) *Job {
	jobLock.Lock()
	defer jobLock.Unlock()
	_id := uuid.New()
	id := _id.String()
	j := &Job{
		Id:     id,
		stream: s,
		Name:   constants.SubjectName + "." + id,
	}
	jobs[j.Id] = j
	return j
}

func (j *Job) StartConsumer(w *websocket.Conn) error {
	return j.stream.StartNewConsumer(j.Id, w)
}
