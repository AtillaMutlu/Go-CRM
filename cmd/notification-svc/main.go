package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Kafka broker adresi
	brokers := []string{"localhost:9092"}
	// Dinlenecek topic
	topic := "notification.command"

	// Kafka reader oluştur
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "notification-svc",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	fmt.Println("Notification servis Kafka'dan mesaj bekliyor...")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		m, err := r.ReadMessage(ctx)
		cancel()
		if err != nil {
			log.Printf("Kafka'dan mesaj okunamadı: %v", err)
			continue
		}
		log.Printf("Yeni bildirim komutu: %s", string(m.Value))
		// Burada email/SMS/WebSocket adapterlerine yönlendirme yapılacak
	}
}
