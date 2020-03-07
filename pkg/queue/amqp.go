package queue

// ----------------------------------------------------------------------
// RabbitMQ
// ----------------------------------------------------------------------
// Link  : http://www.hlzblog.top/
// GITHUB: https://github.com/HaleyLeoZhang
// ----------------------------------------------------------------------
// 文档 https://github.com/streadway/amqp/tree/master/_examples
// ----------------------------------------------------------------------

import (
	"fmt"
	"github.com/streadway/amqp"
	"novels-spider/conf"
	// "time"
)

type AMQP struct {
	// 公共
	Payload []byte
	// 当前驱动配置项
	Exchange string
	Queue    string
}

func (a *AMQP) SetPayload(payload []byte) {
	a.Payload = payload
}

func (a *AMQP) Push() error {

	conn := a.newConnect()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	one.Lock.Lock()
	err = ch.Publish(
		a.Exchange, // exchange
		a.Queue,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        a.Payload,
		})
	one.Lock.Unlock()
	if err != nil {
		return err
	}

	return nil
}

func (a *AMQP) Pull(callback func([]byte) error) error {

	conn := a.newConnect()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	one.Lock.Lock()
	delivery, err := ch.Consume(
		a.Queue, // name
		"",      // consumerTag,
		false,   // noAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // arguments
	)
	one.Lock.Unlock()

	pool := make(chan int, conf.QUEUE_CHANNEL_NUMBER)
	for {
		select {
		case d := <-delivery:
			pool <- 1
			go a.handle(d, callback, pool)
		}
	}
	return nil
}

func (a *AMQP) Close() error {
	err := one.Conn.(*amqp.Connection).Close()
	return err
}

func (a *AMQP) newConnect() *amqp.Connection {
	if nil == one.Conn {
		dial := fmt.Sprintf("amqp://%v:%v@%v:%v/", conf.AMQP["user"], conf.AMQP["password"], conf.AMQP["host"], conf.AMQP["port"])
		var err error
		one.Conn, err = amqp.Dial(dial)
		if err != nil {
			panic([]string{err.Error()})
		}
	}
	return one.Conn.(*amqp.Connection)
}

func (a *AMQP) handle(d amqp.Delivery, callback func([]byte) error, pool chan int) error {

	err := callback(d.Body)
	if err != nil {
		return err
	}
	one.Lock.Lock()
	d.Ack(false)
	one.Lock.Unlock()
	<-pool
	return nil
}