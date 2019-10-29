package client

import (
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/Shopify/sarama"
)

type KafkaSyncProducer struct {
	client sarama.SyncProducer
	debug bool
}

func NewProducer(broker []string, debug bool) *KafkaSyncProducer {
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
		debug: debug,
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
	if sp.debug {
		ilog.Info(fmt.Sprintf("【Kafka Send】 pid: %v, offset: %v", pid, offset))
	}
	return nil
}