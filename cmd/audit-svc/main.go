package main

import (
	"context"
	"fmt"
	"log"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/segmentio/kafka-go"
)

func main() {
	// immudb client başlat
	immu, err := immuclient.NewImmuClient(immuclient.DefaultOptions().WithAddress("localhost").WithPort(3322))
	if err != nil {
		log.Fatalf("immudb bağlantı hatası: %v", err)
	}
	defer immu.Disconnect()

	// Kafka broker adresi
	brokers := []string{"localhost:9092"}
	topic := "audit.raw"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "audit-svc",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	fmt.Println("Audit servis Kafka'dan mesaj bekliyor...")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		m, err := reader.ReadMessage(ctx)
		cancel()
		if err != nil {
			log.Printf("Kafka'dan mesaj okunamadı: %v", err)
			continue
		}
		// Mesajı immudb'ye append et
		_, err = immu.Set(context.Background(), []byte(fmt.Sprintf("audit-%d", m.Offset)), m.Value)
		if err != nil {
			log.Printf("immudb'ye yazılamadı: %v", err)
			continue
		}
		log.Printf("Audit kaydı eklendi: %s", string(m.Value))
	}
}
