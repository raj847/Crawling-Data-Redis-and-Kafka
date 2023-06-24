package main

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func main() {
	// to create topics when auto.create.topics.enable='true'
	_, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "example2", 3)
	if err != nil {
		panic(err.Error())
	}
	
}
