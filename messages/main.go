
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/streadway/amqp"
)

func main() {

	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	// client = initEtcClient()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", user, password, host))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"messages",   // name
		"x-delayed-message", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		amqp.Table {"x-delayed-type": "direct"}, // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	q, err := ch.QueueDeclare(
		"messages-q", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = ch.QueueBind(
		q.Name, // queue name
		q.Name,     // routing key
		"messages", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			msg.Ack(false)
		}
	}()

	handler := func (w http.ResponseWriter, r *http.Request) {
		var delay int
		fmt.Sscanf(r.URL.Query().Get("delay"), "%d", &delay)
		msg := r.URL.Query().Get("msg")
		
		log.Printf("Enqueueing message %s with %dms delay", msg, delay)
		err = ch.Publish(
			"messages",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(msg),
				Headers: amqp.Table{
					"x-delay": delay,
				},
			},
		)
		if err != nil {
			log.Printf("Failed to publish message: %s", err)
		}
	}

    http.HandleFunc("/send", handler)
    log.Print("Listening on port :8888")
    log.Fatal(http.ListenAndServe(":8888", nil))
}
