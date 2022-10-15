package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"quick/conf"
	"quick/pkg/log"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("接收消息: 从话题[ %s ] 发来的内容: %s \n", msg.Topic(), msg.Payload())
}

func New() mqtt.Client {
	mqttConfig := conf.MqttConfig
	//mqttConfig.ClientId = "quick" + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	mqttConfig.ClientId = "quick"
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mqttConfig.Host, mqttConfig.Port))
	opts.SetClientID(mqttConfig.ClientId)
	opts.SetKeepAlive(10 * time.Second)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(3 * time.Second)
	opts.SetCleanSession(false)
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetUsername(mqttConfig.Username)
	opts.SetPassword(mqttConfig.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetWill("offline", "go_mqtt_client offline", 1, false)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Sugar.Errorf("mqtt connect failed:%s", token.Error())
		os.Exit(1)
		return nil
	}
	return c

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	//log.Sugar.Infof("mqtt连接成功...\n")
}

// 连接丢失的回调
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Sugar.Errorf("mqtt连接丢失: %v\n", err.Error())
}
