package jobs

import (
	"bufio"
	"context"
	"errors"
	"log"

	"github.com/Rahul-RB/go-jobqueue/constants"
	"github.com/Rahul-RB/go-jobqueue/utils"
	"github.com/google/uuid"
)

type Job struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	NatsConn *utils.Nats `json:"-"`
}

var jobs map[string]*Job = make(map[string]*Job)

func GetJob(id string) (*Job, error) {
	if j, ok := jobs[id]; ok {
		return j, nil
	}
	return &Job{}, errors.New("Couldn't find job with ID:" + id)
}

func (j *Job) Run() {
	// run job
	// get output of job
	// publish that to NatsConn
	ctx := context.Background()
	cmd, err := utils.RunWithTimeout("/home/rahulrb/go-jobqueue/dummy-job/dummy-job.bin", "-name", j.Name, "-interval", "1s")
	if err != nil {
		log.Fatal("Failed to run command:", err.Error())
		if err := j.NatsConn.Publish(ctx, j.Name, err.Error()); err != nil {
			log.Fatal("Failed to publish to:", j.Name, err.Error())
		}
	}

	scanner := bufio.NewScanner(*cmd.Stdout)
	for scanner.Scan() {
		if err := j.NatsConn.Publish(ctx, j.Name, scanner.Text()); err != nil {
			log.Fatal("Failed to publish to:", j.Name, err.Error())
		}
	}

	cmd.Cmd.Wait()
}

func NewJob(nc *utils.Nats) *Job {
	_id := uuid.New()
	id := _id.String()
	j := &Job{
		Id:       id,
		NatsConn: nc,
		Name:     constants.SubjectName + "." + id,
	}
	jobs[j.Id] = j
	return j
}
