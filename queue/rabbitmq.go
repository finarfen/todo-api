package queue

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	host := getEnv("RABBITMQ_HOST", "localhost")
	port := getEnv("RABBITMQ_PORT", "5672")
	user := getEnv("RABBITMQ_USER", "guest")
	pass := getEnv("RABBITMQ_PASS", "guest")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ недоступен, попытка %d/10...", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Не удалось подключиться к RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Ошибка открытия канала:", err)
	}

	fmt.Println("Подключено к RabbitMQ!")
	return &RabbitMQ{conn: conn, channel: ch}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (r *RabbitMQ) Publish(queueName string, message string) error {
	_, err := r.channel.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	return r.channel.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
