package transfer

import (
	jsoniter "github.com/json-iterator/go"
	"quick/internal/topic"
	"quick/pkg/rabbitmq"
)

type transfer struct {
	metadataWriter rabbitmq.Interface
	dataWriter     rabbitmq.Interface
	eventWriter    rabbitmq.Interface
}
type Writer struct {
}
type Interface interface {
	SendPropertyMetadata(buf []byte)
	SendPropertyData(buf []byte)
	SendPropertyEvent(buf []byte)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func New() Interface {
	t := &transfer{
		metadataWriter: rabbitmq.NewRabbitMQSimple(topic.K_device_metadata_chanl),
		dataWriter:     rabbitmq.NewRabbitMQSimple(topic.K_device_data_chanl),
		eventWriter:    rabbitmq.NewRabbitMQSimple(topic.K_device_event_chanl),
	}

	return t
}
func (t *transfer) SendPropertyMetadata(buf []byte) {
	t.metadataWriter.PushMessage(buf)
}
func (t *transfer) SendPropertyData(buf []byte) {
	t.dataWriter.PushMessage(buf)

}
func (t *transfer) SendPropertyEvent(buf []byte) {
	t.eventWriter.PushMessage(buf)

}
