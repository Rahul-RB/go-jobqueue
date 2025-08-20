package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	name := flag.String("name", "unknown job", "Name for this execution")
	interval := flag.Duration("interval", 1*time.Second, "Interval in seconds")
	flag.Parse()
	for {
		log.Println(*name)
		time.Sleep(*interval)
	}
}
