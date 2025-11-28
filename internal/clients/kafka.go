package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
	kafkaReader *kafka.Reader
)

func InitKafka() {
	// Writer для отправки в chat-in
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("kafka-container:9092"),
		Topic:    "chat-in",
		Balancer: &kafka.LeastBytes{},
	}

	// Reader для чтения из chat-out
	kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka-container:9092"},
		Topic:   "chat-out",
		GroupID: "chat-server-group", // GroupID важен, чтобы читать сообщения только раз
	})
}

func StartKafkaConsumer() {
	log.Println("Kafka Consumer started for topic: chat-out")
	for {
		// Читаем сообщение из Kafka
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading kafka message: %v", err)
			continue
		}

		// Десериализуем (предполагаем, что в chat-out лежит структура KafkaMessage)
		var kMsg KafkaMessage
		if err := json.Unmarshal(m.Value, &kMsg); err != nil {
			log.Printf("Error unmarshalling kafka msg: %v", err)
			continue
		}

		log.Printf("Received from Kafka chat-out for User %s", kMsg.UserID)

		// Ищем соединение пользователя в Hub
		if conn, ok := hub.Get(kMsg.UserID); ok {
			// Отправляем сообщение обратно на фронт
			// В реальности вы можете отправлять kMsg.Body или другую структуру ответа
			err := conn.WriteJSON(kMsg.Body)
			if err != nil {
				log.Printf("Error writing to websocket: %v", err)
				conn.Close()
				hub.Remove(kMsg.UserID)
			}
		} else {
			// Пользователь не подключен к этому инстансу сервера
			// (В микросервисной архитектуре это нормально, он может быть на другой ноде)
			// log.Printf("User %s is not connected locally", kMsg.UserID)
		}
	}
}
