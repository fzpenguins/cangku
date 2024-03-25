package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

var RabbitmqConn *amqp.Connection

func LinkRabbitmq() {

	conn, err := amqp.Dial("amqp://guest:guest@172.23.21.149:5672/")
	if err != nil {
		panic(err)
	}
	RabbitmqConn = conn
	log.Println("connecting to rabbitmq successes")
}
