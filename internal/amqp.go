package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type (
	AMQPService struct {
		client  *amqp091.Connection
		channel *amqp091.Channel
	}
	IAMQPService interface {
		Publish(app string, message Message) (string, error)
		ReplyTo(app string, message MessageReply) (string, error)
		Subscribe(event string) (<-chan amqp091.Delivery, error)
		Close()
	}
)

func MessageToString(message Message) (string, error) {
	jsonBytes, err := json.Marshal(message)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func MessageReplyToString(message MessageReply) (string, error) {
	jsonBytes, err := json.Marshal(message)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func StringToMessage(str string) (Message, error) {
	var message *Message

	err := json.Unmarshal([]byte(str), &message)

	if err != nil {
		return Message{}, err
	}

	return *message, nil
}

func StringToMessageReply(str string) (MessageReply, error) {
	var message *MessageReply

	err := json.Unmarshal([]byte(str), &message)

	if err != nil {
		return MessageReply{}, err
	}

	return *message, nil
}

func (m AMQPService) Publish(app string, message Message) (string, error) {
	q, err := m.channel.QueueDeclare(
		message.Event, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Fatal("failed to declare a queue")
	}

	payload, err := MessageToString(message)

	if err != nil {
		log.Fatal("failed to convert message to string: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = m.channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			AppId:         app,
			ContentType:   "text/plain",
			Body:          []byte(payload),
			CorrelationId: message.CorrelationID, // Set correlation ID in the message properties
		})

	if err != nil {
		return "", fmt.Errorf("failed to publish a message")
	}

	fmt.Print("new message published: ", q.Name+" ", message.CorrelationID, "\n")

	return message.CorrelationID, nil // Return the correlation ID for tracking the response
}

func (m AMQPService) ReplyTo(app string, message MessageReply) (string, error) {
	q, err := m.channel.QueueDeclare(
		message.Event, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Fatal("failed to declare a queue")
	}

	payload, err := MessageReplyToString(message)

	if err != nil {
		log.Fatal("failed to convert message to string: ", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = m.channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			AppId:         app,
			ContentType:   "text/plain",
			Body:          []byte(payload),
			CorrelationId: message.CorrelationID, // Set correlation ID in the message properties
		})

	if err != nil {
		return "", fmt.Errorf("failed to reply a message")
	}

	fmt.Print("new message replied: ", q.Name+" ", message.CorrelationID, "\n")

	return message.CorrelationID, nil // Return the correlation ID for tracking the response
}

func (m AMQPService) Subscribe(event string) (<-chan amqp091.Delivery, error) {
	_, err := m.channel.QueueDeclare(
		event, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %s", err.Error())
	}

	messages, err := m.channel.Consume(
		event, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %s", err.Error())
	}

	return messages, nil
}

func (a AMQPService) Close() {
	err := a.channel.Close()
	if err != nil {
		return
	}
	err = a.client.Close()
	if err != nil {
		return
	}
}

func NewAMQPService() IAMQPService {
	fmt.Println("xx ", os.Getenv("RABBITMQ_DSN"))
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_DSN"))

	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}

	return &AMQPService{
		client:  conn,
		channel: ch,
	}
}
