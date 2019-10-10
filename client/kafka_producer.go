package client

import (
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/Shopify/sarama"
)

type KafkaSyncProducer struct {
	client sarama.SyncProducer
}

func NewProducer(broker []string) *KafkaSyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		ilog.Error("kafka producer create err : ", err)
		return nil
	}

	return &KafkaSyncProducer{
		client: client,
	}
}

func (sp *KafkaSyncProducer) Send(topic string, value string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(value)
	pid, offset, err := sp.client.SendMessage(msg)
	if err != nil {
		return err
	}
	ilog.Info(fmt.Sprintf("kafka send pid:%v offset:%v", pid, offset))
	return nil
}