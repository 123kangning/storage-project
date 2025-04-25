package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"storage/conf"
)

type RabbitMQ struct {
	channel  *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string
}

func init() {
	mq := New(conf.RabbitmqServer)
	//交换机都为扇出，用着简单，直接广播就行
	exchanges := []string{"apiServer", "dataServer"}
	for _, exchange := range exchanges {
		err := mq.channel.ExchangeDeclare(
			exchange,
			"fanout",
			true,  // durable
			false, // auto-deleted
			false, // internal
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Println("ExchangeDeclare err=", err)
			panic(err)
		}
	}
}

func New(s string) *RabbitMQ {
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}

	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}
	//为每个实例生成唯一队列名称
	q, e := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if e != nil {
		panic(e)
	}

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.conn = conn
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	e := q.channel.QueueBind(
		q.Name,   // queue name,之前生成的唯一队列名称
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if e != nil {
		panic(e)
	}
	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	//使用默认直连交换机
	e = q.channel.PublishWithContext(context.Background(), "",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.PublishWithContext(context.Background(), exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, e := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if e != nil {
		panic(e)
	}
	return c
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
	q.conn.Close()
}
