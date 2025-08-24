# go-jobqueue
Distributed job queue written in Go using NATS

## Prerequisites
- Run NATS with jetstream
```
docker run --name nats --network nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222 -j
```
- Create the dummy job binary:
```
cd dummy-job
go build .
```

## Run worker
```bash
go run ./worker
```

## Trigger a job

Each job gets a new Job ID (uuid4)
```bash
curl -X POST http://localhost:3000/v1/job

{
    "id": "db9e2022-3a98-407c-a75f-dec42972c94b",
    "name": "job_output_subject.db9e2022-3a98-407c-a75f-dec42972c94b"
}
```

## Get metadata about the job
```
curl -X GET http://localhost:3000/v1/job/<job-id> && echo
{
    "id": "<job-id>",
    "name": "job_output_subject.<job-id>"
}
```

## Stream job output
This can be done via a websocket client. The easiest way is to:
- open up the browser
- point it towards localhost:3000
- enter the job ID

![stream-output](docs/stream-output.png)

Note: Each new websocket client will read the job output from the very beginning.


## TODO
- filter job output stream by job ID
- keep a limit on number of messages in NATS stream
- allow multiple job types (today we only trigger the dummy job)
