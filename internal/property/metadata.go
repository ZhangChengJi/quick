package property

import (
	"fmt"
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
	Si    int    `json:"si"`
}

type Event struct {
	Iccid string `json:"iccid"`
	Sl    int    `json:"sl"`
	Da    string `json:"da"`
	Le    int    `json:"le"`
}

type Line struct {
	Iccid  string `json:"iccid"`
	Status string `json:"status"`
}
type ApiLine struct {
	DeviceId string `json:"deviceId"`
	Status   int    `json:"status"`
}

const (
	defaultProductId = 0
)

func (m *Metadata) Execute() {
	de, err := queryDevice(m.Iccid) //查询设备
	if err != nil {
		log.Sugar.Errorf("iccid查询失败无法更新slave,%s", err)
		return
	}
	if de.LineStatus == OFFLINE_INT {
		updateDeviceStatus(m.Iccid, ONLINE_INT) //如果设备是离线状态，就改为在线状态
	}
	for _, metadata := range m.List {

		property, err := getProperty(metadata.Pr) //查询属性
		if err != nil {
			return
		}
		if property != nil || metadata.Pr == 0 {
			slave := model.PigDeviceSlave{
				DeviceId:       m.Iccid,
				SlaveName:      "探测器" + strconv.Itoa(metadata.Sl),
				ModbusAddress:  metadata.Sl,
				PropertyId:     metadata.Pr,
				ProductId:      defaultProductId,
				SlaveDesc:      "",
				SlaveStatus:    1,
				LineStatus:     1,
				CreateTime:     time.Now(),
				UpdateTime:     time.Now(),
				InstallDate:    time.Now(),
				CheckStartDate: time.Now(),
				CheckEndDate:   time.Now().AddDate(0, 0, 365),
			}
			createOrUpdateSlave(&slave) //创建或者更新设备的探测器

		} else {
			log.Sugar.Errorf("属性id传递错误%v", err)
			continue
		}
	}
	if len(m.List) > 0 { //如果有探测器
		slaveCount := getSlaveSize(m.Iccid) //查询目前设备的探测器数量
		//如果原来的设备探测器大于新配置的，就进行删除大于出来的设备
		if slaveCount > m.Co { //如果原来的设备探测器大于新配置的，就进行删除大于出来的设备
			as := slaveCount - m.Co     //需要删除的数量
			deleteSlaveMax(m.Iccid, as) //删除探测器
		}
		db.RDB.Del(db.RDB.GetSlaveKey(m.Iccid)) //删除slave缓存
	}
}

const (
	Normal          = iota //正常0
	High                   //高报1
	Low                    //低报2
	Internal               //探测器内部错误3
	Communication          //通讯错误4
	Shield                 //主机屏蔽探测器5
	SlaveHitch             //探测器故障6
	MainHitch              //主电故障
	PrepareHitch           //备电故障
	MainFailback           //主电故障恢复
	PrepareFailback        //备电故障恢复
	failback               //电源全部恢复

	DATA           = "data"
	ALARM          = "alarm"
	HITCH          = "hitch"
	ONLINE         = "online"  //上线
	OFFLINE        = "offline" //下线
	ONLINE_INT     = 1         //上线
	OFFLINE_INT    = 0         //下线
	HostController = 1         //主机
	Detector       = 2         //探测器
)

func (d *Data) Execute() {
	de, err := queryDevice(d.Iccid) //查询设备
	if err != nil {
		return
	}
	if de.LineStatus == OFFLINE_INT { //如果设备是离线状态，就改为在线状态
		updateDeviceStatus(d.Iccid, ONLINE_INT)
	}
	pigSlave, err := querySlave(d.Iccid, d.Sl)
	if err != nil {
		return
	}
	slaveProperty, err := getSlaveProperty(d.Iccid, d.Sl) //查询设备的属性
	if err != nil {
		return
	}
	if slaveProperty != nil { //如果设备属性存在
		var msg *DeviceMsg
		msg = &DeviceMsg{
			Ts: time.Now(),

			DataType:     DATA,
			Level:        d.Le,
			DeviceId:     d.Iccid,
			GroupId:      de.GroupId,
			SlaveId:      d.Sl,
			DeviceType:   Detector,
			SlaveName:    pigSlave.SlaveName,
			Data:         d.Da,
			Unit:         slaveProperty.PropertyUnit,
			PropertyName: slaveProperty.PropertyName,
			Signal:       d.Si,
		}

		queue.Enqueue(msg)                                                                                                                                                 //将数据放入队列
		if d.Le == High || d.Le == Low || d.Le == Internal || d.Le == Communication || d.Le == Shield || d.Le == SlaveHitch || d.Le == MainHitch || d.Le == PrepareHitch { //如果是高报或者低报
			msg = &DeviceMsg{
				Ts: time.Now(),

				GroupId:      de.GroupId,
				DataType:     ALARM,
				Level:        d.Le,
				DeviceId:     d.Iccid,
				SlaveId:      d.Sl,
				DeviceType:   Detector,
				SlaveName:    pigSlave.SlaveName,
				Data:         d.Da,
				Unit:         slaveProperty.PropertyUnit,
				PropertyName: slaveProperty.PropertyName,
				Signal:       d.Si,
			}
			queue.Enqueue(msg)

		}
		Publish(fmt.Sprintf(topic.Device_last, d.Iccid), msg) //将数据发布到mqtt device/last
		if de.GroupId != 0 {
			Publish(fmt.Sprintf(topic.OpenApi_data, strconv.Itoa(de.GroupId), d.Iccid), msg)
		}

	}
}

func (d *Event) Execute() {
	device, err := queryDevice(d.Iccid) //查询设备
	if err != nil {
		return
	}
	if device.LineStatus == OFFLINE_INT { //如果设备是离线状态，就改为在线状态
		updateDeviceStatus(d.Iccid, ONLINE_INT)
	}
	if d.Sl > 0 {
		pigSlave, err := querySlave(d.Iccid, d.Sl)
		if err != nil {
			return
		}
		slaveProperty, err := getSlaveProperty(d.Iccid, d.Sl) //查询设备的属性
		if err != nil {
			return
		}
		if slaveProperty != nil { //如果设备属性存在
			var msg *DeviceMsg
			//         正常			      高             低              探测器内部错误
			if d.Le == Normal || d.Le == High || d.Le == Low || d.Le == Internal || d.Le == Communication || d.Le == Shield || d.Le == SlaveHitch || d.Le == MainHitch || d.Le == PrepareHitch { //如果是高报或者低报或者正常
				if device.GroupId != 0 { //如果设备有分组
					msg = &DeviceMsg{
						Ts: time.Now(),

						DataType:     ALARM,
						Level:        d.Le,
						DeviceId:     d.Iccid,
						SlaveId:      d.Sl,
						GroupId:      device.GroupId,
						DeviceType:   Detector,
						SlaveName:    pigSlave.SlaveName,
						Data:         d.Da,
						Unit:         slaveProperty.PropertyUnit,
						PropertyName: slaveProperty.PropertyName,
						Name:         device.DeviceName,
						Address:      device.DeviceAddress,
					}
					Publish(fmt.Sprintf(topic.Device_event, strconv.Itoa(device.GroupId), d.Iccid), msg) //将数据发布到mqtt device/event
					Publish(fmt.Sprintf(topic.OpenApi_event, strconv.Itoa(device.GroupId), d.Iccid), msg)
					if d.Le == High || d.Le == Low { //如果是高报或者低报
						//有分组的情况下进行发送短信提醒
						//第一次发送过短信需要等待5分钟之后再次发送
						if sendAwait5Second(d.Iccid, d.Sl) {
							//发送电话短信通知
							//topic.Device_notify的+插入d.Ic
							Publish(fmt.Sprintf(topic.Device_notify, d.Iccid), msg) //将数据发布到mqtt device/notify
						}

					}
					if d.Le == High || d.Le == Low || d.Le == Internal || d.Le == Communication || d.Le == Shield || d.Le == SlaveHitch || d.Le == MainHitch || d.Le == PrepareHitch {
						queue.Enqueue(msg) //将数据放入alarm队列
					}

				}

				if d.Le == Normal { //如果为正常
					clearSendAwait(d.Iccid, d.Sl) //清除发送等待
				}

			}
			msg = &DeviceMsg{
				Ts: time.Now(),

				DataType:     DATA,
				Level:        d.Le,
				DeviceId:     d.Iccid,
				GroupId:      device.GroupId,
				DeviceType:   Detector, //代表主机故障还是探测器故障
				SlaveId:      d.Sl,
				SlaveName:    pigSlave.SlaveName,
				Data:         d.Da,
				Unit:         slaveProperty.PropertyUnit,
				PropertyName: slaveProperty.PropertyName,
			}
			queue.Enqueue(msg)                                    //将数据放入正常队列                              //将数据放入队列
			Publish(fmt.Sprintf(topic.Device_last, d.Iccid), msg) //将数据发布到mqtt device/last
			Publish(fmt.Sprintf(topic.OpenApi_data, strconv.Itoa(device.GroupId), d.Iccid), msg)
		}
	}
	if d.Sl == 0 { //0代表是主机
		var msg *DeviceMsg

		updateDeviceHitch(d.Iccid, d.Le) //修改7主电8备电故障+9主电10备电恢复 11全部恢复
		if d.Le == MainHitch || d.Le == PrepareHitch {
			if device.GroupId != 0 { //如果设备有分组
				msg = &DeviceMsg{
					Ts: time.Now(),

					DataType:   ALARM,
					Level:      d.Le,
					DeviceId:   d.Iccid,
					GroupId:    device.GroupId,
					DeviceType: HostController, //代表主机故障还是探测器故障
					Name:       device.DeviceName,
					Address:    device.DeviceAddress,
				}
				queue.Enqueue(msg)
				Publish(fmt.Sprintf(topic.Device_notify, d.Iccid), msg)
				Publish(fmt.Sprintf(topic.Device_event, strconv.Itoa(device.GroupId), d.Iccid), msg) //将数据发布到mqtt device/event
				Publish(fmt.Sprintf(topic.OpenApi_event, strconv.Itoa(device.GroupId), d.Iccid), msg)

			}
		}
		msg = &DeviceMsg{
			Ts:         time.Now(),
			DataType:   DATA,
			Level:      d.Le,
			DeviceId:   d.Iccid,
			GroupId:    device.GroupId,
			DeviceType: HostController, //代表主机故障还是探测器故障
		}
		queue.Enqueue(msg)                                    //将数据放入正常队列                              //将数据放入队列
		Publish(fmt.Sprintf(topic.Device_last, d.Iccid), msg) //将数据发布到mqtt device/last
		Publish(fmt.Sprintf(topic.OpenApi_data, strconv.Itoa(device.GroupId), d.Iccid), msg)
	}

}
func (l *Line) Execute() {

	if l.Status == ONLINE {
		updateDeviceStatus(l.Iccid, 1)
	} else {
		updateDeviceStatus(l.Iccid, 0)
	}
}
