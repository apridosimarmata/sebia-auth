package infrastructure

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	nsq "github.com/nsqio/go-nsq"
)

type MessagingProducer interface {
	PublishMessage(ctx context.Context, topic string, channel string, message interface{}) (err error)
}

type messagingProducer struct {
	nsqProducer *nsq.Producer
}

func NewMessagingProducer() MessagingProducer {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Panic(err)
	}

	return &messagingProducer{
		nsqProducer: p,
	}
}

func (producer *messagingProducer) PublishMessage(ctx context.Context, topic string, channel string, message interface{}) (err error) {

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = producer.nsqProducer.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}

// consumer
type MessagingConsumerInterface interface {
	ConsumeMessage(message *nsq.Message) (err error)
}

type RegisterListenersParam struct {
	Topic    string
	Channel  string
	Listener MessagingConsumerInterface
}

func RegisterConsumers(params []RegisterListenersParam) {
	wg := &sync.WaitGroup{}
	wg.Add(len(params))

	for _, param := range params {
		copyParam := param
		go func(param RegisterListenersParam) {
			decodeConfig := nsq.NewConfig()
			c, err := nsq.NewConsumer(param.Topic, param.Channel, decodeConfig)
			if err != nil {
				log.Panic("Could not create consumer")
			}

			c.AddHandler(nsq.HandlerFunc(param.Listener.ConsumeMessage))

			err = c.ConnectToNSQD("127.0.0.1:4150")
			if err != nil {
				log.Panic("Could not connect")
			}
		}(copyParam)
	}

	wg.Wait()
}
