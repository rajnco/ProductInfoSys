package rmqsender

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type Sender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
}

func Connect(queueName string) *Sender {
	//connection, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	//connection, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Panicf("failed to connect RabbitMQ : %+v ", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Panicf("failed to get channel : %+v ", err)
	}

	//queue, err := channel.QueueDeclare("Produced", false, false, false, false, nil)
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Panicf("failed to get queue : %+v ", err)
	}

	return &Sender{connection: connection, channel: channel, queue: &queue}
}

func (s *Sender) SendMessage(ctx context.Context, content string) error {
	message := amqp.Publishing{ContentType: "text/plain", Body: []byte(content)}
	return s.channel.PublishWithContext(ctx, "", s.queue.Name, false, false, message)
}

func (s *Sender) Close() {
	s.channel.Close()
	s.connection.Close()
}
