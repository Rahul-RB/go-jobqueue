# go-jobqueue
Distributed job queue written in Go using NATS

## Run worker
```bash
go run ./worker
```

## Send a curl request
```bash
curl -vvv -X POST http://localhost:3000/v1/job && echo
```

## Check NATS consumer
```bash
nats consumer next job_output_stream job_output_stream_consumer
```
