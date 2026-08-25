package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/events"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/kafka/consumer"
)

// Обработчик для лидерборда
func handleActivity(msg []byte) error {
	var act events.Activity
	if err := json.Unmarshal(msg, &act); err != nil {
		return err
	}

	switch act.Type {
	case "like":
		updateScore(act.UserID, 1)
	case "comment":
		updateScore(act.UserID, 5)
	case "post":
		updateScore(act.UserID, 10)
	}
	return nil
}

func updateScore(userID uuid.UUID, points int) {
	fmt.Println(userID, points)
}

func main() {
	brokers := []string{"kafka:9092"}
	consumer, err := consumer.NewMyConsumer(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// Передаем топик + функцию-обработчик
	if err := consumer.Consume("activities", handleActivity); err != nil {
		log.Fatal(err) // в отдельной горутине запускается
	}
}
