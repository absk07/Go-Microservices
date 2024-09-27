package events

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Method       string `json:"method"`
	Path         string `json:"path"`
	RemoteAddr   string `json:"remote_addr"`
	ResponseTime string `json:"response_time"`
	StartTime    string `json:"start_time"`
	StatusCode   string `json:"status_code"`
	Severity     string `json:"severity"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return nil
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return nil
	}

	for _, topic := range topics {
		err := ch.QueueBind(
			q.Name,
			topic,
			"logs",
			false,
			nil,
		)
		if err != nil {
			return nil
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			var payload Payload
			_ = json.Unmarshal(msg.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for messages [Exchange, Queue] [logs, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	// switch payload.Name {
	// case "log", "event":
	// 	err := logEvent(payload)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// case "auth":
	// 	// authenticate
	// default:
	// 	err := logEvent(payload)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	err := logEvent(payload)
	if err != nil {
		log.Println(err)
	}
}

func logEvent(payload Payload) error {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "http://logger-service:7070/log", bytes.NewBuffer(jsonData))

	if err != nil {
		// log.Println("http req err", err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		// log.Println("http res err", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
