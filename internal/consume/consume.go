package consume

import (
	"encoding/json"
	"quick/internal/property"
	"quick/internal/topic"
	"quick/pkg/log"
	"quick/pkg/rabbitmq"
)

type consume struct {
	metadataRead rabbitmq.Interface
	dataRead     rabbitmq.Interface
	eventRead    rabbitmq.Interface
	lineRead     rabbitmq.Interface
}
type Interface interface {
	receivedPropertyMetadata()
	receivedPropertyData()
	receivedPropertyEvent()
}

func Start() {
	t := &consume{
		metadataRead: rabbitmq.NewRabbitMQSimple(topic.K_device_metadata_chanl),
		dataRead:     rabbitmq.NewRabbitMQSimple(topic.K_device_data_chanl),
		eventRead:    rabbitmq.NewRabbitMQSimple(topic.K_device_event_chanl),
		lineRead:     rabbitmq.NewRabbitMQSimple(topic.K_device_line),
	}
	go t.metadataRead.ConsumeSimple(t.receivedPropertyMetadata)
	go t.eventRead.ConsumeSimple(t.receivedPropertyEvent)
	go t.dataRead.ConsumeSimple(t.receivedPropertyData)
	go t.lineRead.ConsumeSimple(t.receivedDeviceLine)
}

func (c *consume) receivedPropertyMetadata(bytes []byte) {
	var me *property.Metadata
	if err := c.format(bytes, &me); err == nil {
		me.Execute()
	}
}
func (c *consume) receivedPropertyData(bytes []byte) {
	var me *property.Data
	if err := c.format(bytes, &me); err == nil {
		me.Execute()
	}
}
func (c *consume) receivedPropertyEvent(bytes []byte) {
	var me *property.Event
	if err := c.format(bytes, &me); err == nil {
		me.Execute()
	}
}
func (c *consume) receivedDeviceLine(bytes []byte) {
	var me *property.Line
	if err := c.format(bytes, &me); err == nil {
		me.Execute()
	}
}
func (c *consume) format(msg []byte, i interface{}) error {
	err := json.Unmarshal(msg, i)
	if err != nil {
		log.Sugar.Errorf("格式序列化错误%s", err)
		return err
	}
	return nil
}
