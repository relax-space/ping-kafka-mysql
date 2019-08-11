package eventconsume

import (
	"context"

	"github.com/Shopify/sarama"

	"github.com/pangpanglabs/goutils/kafka"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type EventConsumer struct {
	brokers []string
	topic   string
	groupId string
}

type HandlerFunc func(ConsumeContext) error

type ConsumeContext interface {
	Bind(v interface{}) error
	Context() context.Context
	Status() string
}

func NewEventConsumer(groupId string, brokers []string, topic string) *EventConsumer {
	c := EventConsumer{
		brokers: brokers,
		topic:   topic,
		groupId: groupId,
	}

	return &c
}

func (c *EventConsumer) Handle(f HandlerFunc) error {
	//sarama.Logger = log.New(os.Stderr, "[Sarama] ", log.LstdFlags)
	consumer, err := kafka.NewConsumerGroup(c.groupId, c.brokers, c.topic, func(config *sarama.Config) {
		//config.Metadata.Retry.Backoff = 1 * time.Second
		//config.Metadata.Retry.Max = 300
	})
	if err != nil {
		return err
	}

	messages, err := consumer.Messages()
	if err != nil {
		return err
	}

	go func() {
		for m := range messages {
			status := jsoniter.Get(m.Value, "status").ToString()
			logEntry := logrus.WithFields(logrus.Fields{
				"Offset":    m.Offset,
				"Partition": m.Partition,
				"Topic":     m.Topic,
				"Status":    status,
			})

			handler := func(ctx context.Context) error {
				c := consumeContext{
					value:  m.Value,
					ctx:    ctx,
					status: status,
				}
				return f(c)
			}

			if err := handler(context.Background()); err != nil {
				logEntry.WithError(err).Error("Fail to consume event")
				continue
			}

			logEntry.Info("Success to consume event")
		}
	}()

	return nil
}

type consumeContext struct {
	value  []byte
	ctx    context.Context
	status string
}

func (c consumeContext) Bind(v interface{}) error {
	jsoniter.Get(c.value, "payload").ToVal(v)
	return nil
}

func (c consumeContext) Context() context.Context {
	return c.ctx
}

func (c consumeContext) Status() string {
	return c.status
}
