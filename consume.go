package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RecivedReviewToRabbitmq(conn *amqp.Connection) {
	channel, err := conn.Channel()
	if err != nil {
		log.Println("failed to create channel", err)
	}

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

		}
	}()

	<-forever

}

// func sendEmail(email, subject, body string) error {
// 	from := os.Getenv("GMAIL_ACCOUNT")
// 	password := os.Getenv("GMAIL_SECRET")
// 	host := "smtp.gmail.com"
// 	port := 587

// 	// Connect to the SMTP server.
// 	auth := smtp.PlainAuth("", from, password, host)
// 	to := []string{email}
// 	msg := []byte("To: " + email + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		"Comment:" + body + "\r\n")
// 	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, to, msg)
// 	if err != nil {
// 		return nil
// 	}
// 	log.Println("email is successfully sent to " + email)
// 	return err
// }

// // sendEmail("shuaibuabdulkadir222@gmail.com", "Shuaibu comments on your product", "shuayb")

// func sendEmail(email, subject, body string) {
// 	from := os.Getenv("GMAIL_ACCOUNT")
// 	password := os.Getenv("GMAIL_SECRET")
// 	host := "smtp.gmail.com"
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", from)
// 	m.SetHeader("To", from)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", "Hello <b>Shuaibu</b> and <i>Abdulkadir</i>!")

// 	d := gomail.NewDialer(host, 587, from, password)

// 	if err := d.DialAndSend(m); err != nil {
// 		log.Println(err)
// 		return
// 	}
// }
