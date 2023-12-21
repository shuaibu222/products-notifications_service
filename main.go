package main

import (
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	log.Println("Notification service running...")

	// initialize the connection of rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		os.Exit(1)
	}

	defer conn.Close()

	go RecivedFromRabbitmq("users", conn)
	go RecivedFromRabbitmq("reviews", conn)
	go RecivedFromRabbitmq("products", conn)

	// solved
	select {}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
