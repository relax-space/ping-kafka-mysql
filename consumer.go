package main

import (
	"fmt"
	"nomni/utils/eventconsume"
)

type Config struct {
	Brokers []string
	Topic   string
}

type Fruit struct {
	Name string `json:"name"`
}

func Consume(serviceName string, kafkaConfig Config,
	f func(eventconsume.ConsumeContext) error) error {
	return eventconsume.NewEventConsumer(
		serviceName,
		kafkaConfig.Brokers,
		kafkaConfig.Topic,
	).Handle(f)
}

func EventFruit(c eventconsume.ConsumeContext) error {
	var fruit Fruit
	if err := c.Bind(&fruit); err != nil {
		return err
	}
	if c.Status() == "FruitCreated" {
		fmt.Printf("accepted message:%+v \n", fruit)
	}
	return nil
}
