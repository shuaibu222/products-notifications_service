package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (app *Config) RecivedReviewToRabbitmq(conn *amqp.Connection) {
	channel, err := conn.Channel()
	if err != nil {
		log.Println("failed to create channel", err)
	}

	defer conn.Close()

	defer channel.Close()

	q, err := channel.QueueDeclare(
		"reviews", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Println("failed to declare a queue", err)
	}

	// We use exchange when we want producer to send to different queues without interacting directly with queue
	// err = channel.ExchangeDeclare(
	// 	"reviews_exchange", // Exchange name
	// 	"direct",           // Exchange type (or 'topic' if needed)
	// 	true,               // Durable
	// 	false,              // Auto-deleted
	// 	false,              // Internal
	// 	false,              // No-wait
	// 	nil,                // Arguments
	// )
	// if err != nil {
	// 	log.Println(err)
	// }

	// Bind the queue to the exchange to let them know each other. with that we can have as many queues as we want to the same exchange
	// err = channel.QueueBind(
	// 	q.Name,             // Queue name
	// 	"reviews",          // Routing key
	// 	"reviews_exchange", // Exchange
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	log.Println(err)
	// }

	msgs, err := channel.Consume(
		q.Name, // routing key
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // mandatory
		false,  // immediate
		nil,
	)

	if err != nil {
		log.Println("failed to consume", err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Recived a message: %s", d.Body)
			err := sendEmail("shuaibuabdulkadir222@gmail.com", "Shuaibu comments on your product", string(d.Body))
			if err != nil {
				log.Println("Failed to send email", err)
			}
		}
	}()

	<-forever
}

func sendEmail(email, subject, body string) error {
	from := os.Getenv("GMAIL_ACCOUNT")
	password := os.Getenv("GMAIL_SECRET")
	host := "smtp.gmail.com"
	port := 587

	// Connect to the SMTP server.
	auth := smtp.PlainAuth("", from, password, host)
	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		"Comment:" + body + "\r\n")
	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, to, msg)
	if err != nil {
		return nil
	}
	log.Println("email is successfully sent to " + email)
	return err
}
