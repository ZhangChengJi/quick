package property

import (
	"quick/internal/topic"
	"quick/manager/model"
	"quick/pkg/log"
	"strconv"
	"time"
)

type Metadata struct {
	Iccid string           `json:"iccid"`
	Co    int              `json:"co"`
	List  []*SlaveMetadata `json:"list"`
}
type SlaveMetadata struct {
	Sl int `json:"sl"`
	Pr int `json:"pr"`
}

type Data struct {
	Iccid string `json:"iccid"`
	Sl    int    `json:"sl"`
	Da    string `json:"da"`
	Le    int    `json:"le"`
}
type Event struct {
	Iccid string `json:"iccid"`
	Sl    int    `json:"sl"`
	Da    string `json:"da"`
	Le    int    `json:"le"`
}

func (m *Metadata) Execute() {
	_, err := queryDevice(m.Iccid)
	if err != nil {
		log.Sugar.Errorf("iccid查询失败无法更新slave,%s", err)
		return
	}
	for _, metadata := range m.List {
		slave := model.PigDeviceSlave{
			DeviceId:      m.Iccid,
			SlaveName:     "探头" + strconv.Itoa(metadata.Sl),
			ModbusAddress: metadata.Sl,
			PropertyId:    metadata.Pr,
			SlaveDesc:     "",
			SlaveStatus:   1,
			LineStatus:    1,
		}
		createOrUpdateSlave(&slave)
	}
}

const (
	Normal        = iota //正常
	High                 //高报
	Low                  //低报
	Internal             //探测器内部错误
	Communication        //通讯错误
	Shield               //主机屏蔽探测器
	SlaveHitch           //探头故障
	DATA          = "data"
	ALARM         = "alarm"
	HITCH         = "hitch"
)

func (d *Data) Execute() {

	slaveProperty, err := getSlaveProperty(d.Iccid, d.Sl)
	if err != nil {
		return
	}
	msg := &DeviceMsg{
		Ts:           time.Now(),
		DataType:     DATA,
		Level:        d.Le,
		DeviceId:     d.Iccid,
		SlaveId:      d.Sl,
		SlaveName:    "探头" + strconv.Itoa(d.Sl),
		Data:         d.Da,
		Unit:         slaveProperty.PropertyUnit,
		PropertyName: slaveProperty.PropertyName,
	}
	queue.Enqueue(msg)
	if d.Le == High || d.Le == Low {
		msg := &DeviceMsg{
			Ts:           time.Now(),
			DataType:     ALARM,
			Level:        d.Le,
			DeviceId:     d.Iccid,
			SlaveId:      d.Sl,
			SlaveName:    "探头" + strconv.Itoa(d.Sl),
			Data:         d.Da,
			Unit:         slaveProperty.PropertyUnit,
			PropertyName: slaveProperty.PropertyName,
		}
		queue.Enqueue(msg)
	}

}

func (d *Event) Execute() {
	device, err := queryDevice(d.Iccid)
	if err != nil {
		return
	}
	slaveProperty, err := getSlaveProperty(d.Iccid, d.Sl)
	if err != nil {
		return
	}

	msg := &DeviceMsg{
		Ts:           time.Now(),
		DataType:     ALARM,
		Level:        d.Le,
		DeviceId:     d.Iccid,
		SlaveId:      d.Sl,
		SlaveName:    "探头" + strconv.Itoa(d.Sl),
		Data:         d.Da,
		Unit:         slaveProperty.PropertyUnit,
		PropertyName: slaveProperty.PropertyName,
		Name:         device.DeviceName,
		Address:      device.DeviceAddress,
	}
	//有分组的情况下进行发送短信提醒
	if device.GroupId != 0 {
		//第一次发送过短信需要等待5分钟之后再次发送
		if sendAwait5Second(d.Iccid, d.Sl) {
			//发送电话短信通知
			Publish(string(append([]byte(topic.Device_notify), d.Iccid...)), msg)
		}
	}
	queue.Enqueue(msg)

	msg1 := &DeviceMsg{
		Ts:           time.Now(),
		DataType:     DATA,
		Level:        d.Le,
		DeviceId:     d.Iccid,
		SlaveId:      d.Sl,
		SlaveName:    "探头" + strconv.Itoa(d.Sl),
		Data:         d.Da,
		Unit:         slaveProperty.PropertyUnit,
		PropertyName: slaveProperty.PropertyName,
	}
	queue.Enqueue(msg1)
}
