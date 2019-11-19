package kafka

import (
	"context"
	"fmt"
	"github.com/CalvinDjy/iteaGo/constant"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/ioc/iface"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"strings"
)

const (
	KEY_KEY 	= "Key"
	HANDLER_KEY = "Handler"
)

type KafkaConsumer struct {
	Ctx             context.Context
	Ioc 			iface.IIoc
	Name 			string
	Brokers			string
	Topic			string
	Group			string
	Processor		[]interface{}
	consumer		*cluster.Consumer
	handler			map[string][]IHandler
	debug 			bool
}

func (kc *KafkaConsumer) Execute() {

	if d, ok := kc.Ctx.Value(constant.DEBUG).(bool); ok {
		kc.debug = d
	}

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Group.Mode = cluster.ConsumerModePartitions
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	
	// init consumer
	if kc.Brokers == "" {
		panic("kafka broker can not be empty")
	}
	brokers := strings.Split(kc.Brokers, ",")
	
	if kc.Topic == "" {
		panic("kafka topic can not be empty")
	}
	topics := []string{kc.Topic}
	
	if kc.Group == "" {
		kc.Group = fmt.Sprintf("%s%s%d", kc.Topic, "_group_", 1)
	}
	
	var err error
	kc.consumer, err = cluster.NewConsumer(brokers, kc.Group, topics, config)
	if err != nil {
		panic(err)
	}
	
	// consume errors
	go func() {
		for err := range kc.consumer.Errors() {
			ilog.Error("consumer error: ", err.Error())
		}
	}()
	
	// consume notifications
	go func() {
		for ntf := range kc.consumer.Notifications() {
			ilog.Info(fmt.Sprintf("consumer rebalanced: %+v", ntf))
		}
	}()

	kc.initHandler()
	ilog.Info(fmt.Sprintf("=== 【Kafka】Consumer [%s] start [Topic : %s, Group : %s] ===", kc.Name, kc.Topic, kc.Group))
	kc.start()
}

func (kc *KafkaConsumer) initHandler() {
	kc.handler = map[string][]IHandler{}
	for _, v := range kc.Processor {
		
		if i, ok := v.(string); ok {
			kc.appendHandler("", i)
			continue
		}

		if i, ok := v.(map[interface{}]interface{}); ok {
			k := ""
			if _, ok := i[KEY_KEY]; ok {
				k = i[KEY_KEY].(string)
			}
			kc.appendHandler(k, i[HANDLER_KEY].(string))
		}

	}
}

func (kc *KafkaConsumer) appendHandler(k string, i string) {
	h := kc.Ioc.InsByName(i)
	if h == nil {
		ilog.Error(fmt.Sprintf("consumer [%s] is nil, please check out if [%s] is registed", i, i))
		return
	}

	if v, ok := h.(IHandler); ok {
		if _, ok := kc.handler[k]; !ok {
			kc.handler[k] = []IHandler{}
		}
		kc.handler[k] = append(kc.handler[k], v)
		return
	}

	ilog.Error(fmt.Sprintf("consumer [%s] is not impliment of kafka.IHandler", i))
}

func (kc *KafkaConsumer) start () {
	stop := make(chan bool, 1)
	go kc.stop(stop)

	// consume messages, watch signals
	for {
		select {
		case part, ok := <-kc.consumer.Partitions():
			if !ok {
				return
			}

			// start a separate goroutine to consume messages
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					if kc.debug {
						ilog.Info(fmt.Sprintf("【Kafka Receive】 partition: %d, key: %s, topic: %s, value: %s", msg.Partition, msg.Key, msg.Topic, msg.Value))
					}
					kc.deal(msg)
				}
			}(part)
		case <-stop:
			kc.consumer.Close()
			return
		}
	}
}

func (kc *KafkaConsumer) deal(msg *sarama.ConsumerMessage) {
	
	var handlerList []IHandler

	msgKey := string(msg.Key)

	if l, ok := kc.handler[msgKey]; ok {
		handlerList = append(handlerList, l...)
	}

	if msgKey != "" {
		if l, ok := kc.handler[""]; ok {
			handlerList = append(handlerList, l...)
		}
	}

	if len(handlerList) == 0 {
		ilog.Error(fmt.Sprintf("message key [%s] has not matched handler", msg.Key))
	} else {
		for _, h := range handlerList {
			err := h.DealMessage(msg.Topic, msg.Partition, msg.Value)
			if err == nil {
				kc.consumer.MarkOffset(msg, "") // mark message as processed
			}
		}
	}
}

//KafkaConsumer stop
func (kc *KafkaConsumer) stop(stop chan bool) {
	for {
		select {
		case <-	kc.Ctx.Done():
			ilog.Info("kafka consumer stop ...")
			stop <- true
			ilog.Info("kafka consumer stop success")
			return
		}
	}
}