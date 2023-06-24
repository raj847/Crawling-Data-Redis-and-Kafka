// to produce messages
package main

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// topic := "example"
	// partition := 0

	// conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	// if err != nil {
	// 	log.Fatal("failed to dial leader:", err)
	// }

	// conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	// _, err = conn.WriteMessages(
	// 	kafka.Message{Value: []byte("saya")},
	// 	kafka.Message{Value: []byte("adalah")},
	// 	kafka.Message{Value: []byte("manusia")},
	// )
	// if err != nil {
	// 	log.Fatal("failed to write messages:", err)
	// }

	// if err := conn.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }

	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "example2",
		Balancer: &kafka.RoundRobin{},
	}

	for {
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("d"),
				Value: []byte("ddd"),
			},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	//hrrmjatfwuuemrlz
}
