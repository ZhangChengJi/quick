package property

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"quick/internal/topic"
	"quick/internal/transfer"
	"quick/manager/database"
	mq "quick/pkg/mqtt"
	queue2 "quick/pkg/queue"

	"strings"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type property struct {
	mqtt     mqtt.Client
	transfer transfer.Interface
	zq1      chan *Metadata
	zq2      chan *Data
	zq3      chan *Event
}

type Interface interface {
	listenMqtt()
	propertyMetaHandler(client mqtt.Client, msg mqtt.Message)
	propertyHandler(client mqtt.Client, msg mqtt.Message)
	eventHandler(client mqtt.Client, msg mqtt.Message)
	connectHandler(client mqtt.Client, msg mqtt.Message)
	disconnectHandler(client mqtt.Client, msg mqtt.Message)
}

func Start() {
	pr := &property{
		mqtt:     mq.New(),
		transfer: transfer.New(),
	}
	db = database.New()
	queue = queue2.NewWithOption(queue2.DefaultOption())
	mqt = pr.mqtt
	time.Sleep(300 * time.Microsecond)
	pr.run()

}
func (p *property) run() {
	p.zq1 = make(chan *Metadata, 1024)
	p.zq2 = make(chan *Data, 1024)
	p.zq3 = make(chan *Event, 1024)
	go p.zqRead()
	go batch(queue, 4, db.TDB)
	p.mqtt.Subscribe(topic.Property_config_post_topic, 1, p.propertyMetaHandler)
	p.mqtt.Subscribe(topic.Property_post_topic, 1, p.propertyHandler)
	p.mqtt.Subscribe(topic.Event_post_topic, 1, p.eventHandler)
}
func (p *property) propertyMetaHandler(client mqtt.Client, msg mqtt.Message) {
	if iccid, ok := p.getIccid(msg.Topic()); ok {
		var metadata *Metadata
		if err := p.format(msg.Payload(), &metadata); err == nil {
			metadata.Iccid = iccid
			p.zq1 <- metadata
		}
	}
}

func (p *property) propertyHandler(client mqtt.Client, msg mqtt.Message) {
	if iccid, ok := p.getIccid(msg.Topic()); ok {
		var data []*Data
		if err := p.format(msg.Payload(), &data); err == nil {
			for _, da := range data {
				da.Iccid = iccid
				p.zq2 <- da
			}
		}
	}
}
func (p *property) eventHandler(client mqtt.Client, msg mqtt.Message) {
	if iccid, ok := p.getIccid(msg.Topic()); ok {

		var event *Event
		if err := p.format(msg.Payload(), &event); err == nil {
			event.Iccid = iccid
			p.zq3 <- event
		}
	}
}

func (p *property) connectHandler(client mqtt.Client, msg mqtt.Message) {

}
func (p *property) disconnectHandler(client mqtt.Client, msg mqtt.Message) {

}
func (p *property) zqRead() {
	for {
		select {
		case entry := <-p.zq1:
			var buf []byte
			p.tobuf(entry, &buf)
			p.transfer.SendPropertyMetadata(buf)
		case entry := <-p.zq2:
			var buf []byte
			p.tobuf(&entry, &buf)
			p.transfer.SendPropertyData(buf)
		case entry := <-p.zq3:
			var buf []byte
			p.tobuf(&entry, &buf)
			p.transfer.SendPropertyEvent(buf)
		}

	}
}

func (p *property) format(msg []byte, t interface{}) error {
	err := json.Unmarshal(msg, t)
	if err != nil {
		return err
	}
	return nil
}
func (p *property) tobuf(entry interface{}, buf *[]byte) error {
	jsonBuf, err := json.Marshal(entry)
	*buf = append(*buf, jsonBuf...)
	return err
}

func (p *property) getIccid(topic string) (string, bool) {
	split := strings.Split(topic, "/")
	if len(split) > 2 {
		return split[1], true
	} else {
		return "", false
	}
}
