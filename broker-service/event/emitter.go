package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println("emitter.go:setup(). Error occurred during get channel from emitter...")
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		log.Println("emitter.go:Push(). Error occurred during get channel from emitter...")
		return err
	}
	defer channel.Close()

	log.Println("Pushing to the channel")
	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		log.Println("emitter.go:Push(). Error occurred during publishing the event to the channel...")
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		if err != nil {
			log.Println("emitter.go:NewEventEmitter(). Error occurred during Emitter setup...")
			return Emitter{}, err
		}
	}

	return emitter, nil
}
