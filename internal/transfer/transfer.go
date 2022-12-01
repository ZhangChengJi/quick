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
	deviceLine     rabbitmq.Interface
}
type Writer struct {
}
type Interface interface {
	SendPropertyMetadata(buf []byte)
	SendProperty(buf []byte)
	SendEvent(buf []byte)
	SendLine(buf []byte)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func New() Interface {
	t := &transfer{
		metadataWriter: rabbitmq.NewRabbitMQSimple(topic.K_device_metadata_chanl),
		dataWriter:     rabbitmq.NewRabbitMQSimple(topic.K_device_data_chanl),
		eventWriter:    rabbitmq.NewRabbitMQSimple(topic.K_device_event_chanl),
		deviceLine:     rabbitmq.NewRabbitMQSimple(topic.K_device_line),
	}

	return t
}

func (t *transfer) SendPropertyMetadata(buf []byte) {
	t.metadataWriter.PushMessage(buf)
}
func (t *transfer) SendProperty(buf []byte) {
	t.dataWriter.PushMessage(buf)
}
func (t *transfer) SendEvent(buf []byte) {
	t.eventWriter.PushMessage(buf)
}
func (t *transfer) SendLine(buf []byte) {
	t.deviceLine.PushMessage(buf)
}
