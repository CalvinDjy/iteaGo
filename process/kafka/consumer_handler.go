package kafka

type IHandler interface {
	DealMessage(topic string, partition int32, value []byte) error
}