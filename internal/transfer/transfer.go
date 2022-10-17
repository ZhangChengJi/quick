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
	SendPropertyData(buf []byte)
	SendPropertyEvent(buf []byte)
	SendDeviceLine(buf []byte)
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
func (t *transfer) SendPropertyData(buf []byte) {
	t.dataWriter.PushMessage(buf)
}
func (t *transfer) SendPropertyEvent(buf []byte) {
	t.eventWriter.PushMessage(buf)
}
func (t *transfer) SendDeviceLine(buf []byte) {
	t.deviceLine.PushMessage(buf)
}
