package property

import (
	"fmt"
	"quick/internal/topic"
	"quick/manager/model"
	"strconv"
	"time"
)

func queryDevice(iccid string) (*model.PigDevice, error) {
	var pigDevice *model.PigDevice
	if err := db.RDB.Get(db.RDB.GetDeviceKey(iccid), &pigDevice); err == nil {
		return pigDevice, nil

	}
	err := db.DB.Where(&model.PigDevice{Id: iccid}).First(&pigDevice).Error
	if err == nil {
		marshal, err := json.Marshal(&pigDevice)
		if err != nil {
			return nil, err
		}
		db.RDB.Set(db.RDB.GetDeviceKey(iccid), string(marshal), -1)
		return pigDevice, nil
	}
	return pigDevice, err
}
func updateDeviceStatus(iccid string, status int) {
	device, err := queryDevice(iccid)
	if err != nil {
		return
	}
	if device.LineStatus == status {
		return
	}
	if device.GroupId != 0 {
		Publish(fmt.Sprintf(topic.Device_line, strconv.Itoa(device.GroupId), iccid), status)
	}
	var pigDevice *model.PigDevice
	err = db.DB.Model(&pigDevice).Where("id=?", iccid).Update("line_status", status).Error
	var pigDeviceSlave model.PigDeviceSlave
	err = db.DB.Model(&pigDeviceSlave).Where("device_id=? ", iccid).Update("line_status", status).Error
	if err == nil {
		err := db.DB.Where(&model.PigDevice{Id: iccid}).First(&pigDevice).Error
		if err == nil {
			marshal, err := json.Marshal(&pigDevice)
			if err != nil {
				return
			}
			db.RDB.Set(db.RDB.GetDeviceKey(iccid), string(marshal), -1)
		}
	}

}
func createOrUpdateSlave(slave *model.PigDeviceSlave) {
	var pigDeviceSlave *model.PigDeviceSlave
	err := db.DB.Where(&model.PigDeviceSlave{
		DeviceId:      slave.DeviceId,
		ModbusAddress: slave.ModbusAddress,
	}).First(&pigDeviceSlave).Error
	if err != nil { //不存在
		db.DB.Create(&slave) //创建
	} else {
		if slave.PropertyId != 0 && slave.PropertyId != pigDeviceSlave.PropertyId { //如果新传递的属性不为0，并且不等于原来的属性
			pigDeviceSlave.PropertyId = slave.PropertyId //更新属性
			db.DB.Model(&pigDeviceSlave).UpdateColumns(&model.PigDeviceSlave{PropertyId: slave.PropertyId})

		}
	}

}
func getSlaveSize(iccid string) int {
	var count int64
	db.DB.Debug().Model(&model.PigDeviceSlave{}).Where("device_id=?", iccid).Count(&count)
	strInt64 := strconv.FormatInt(count, 10)
	id16, _ := strconv.Atoi(strInt64)
	return id16
}
func deleteSlaveMax(iccid string, size int) {
	if len(iccid) > 0 {
		db.DB.Debug().Where("device_id=?", iccid).Order("modbus_address DESC").Limit(size).Delete(&model.PigDeviceSlave{})
	}
}

func getSlaveProperty(iccid string, slaveId int) (*model.PigProperty, error) {
	var slave model.PigDeviceSlave
	err := db.RDB.HGet(db.RDB.GetSlaveKey(iccid), strconv.Itoa(slaveId), &slave)
	if err == nil {
		return getProperty(slave.PropertyId)
	}
	err = db.DB.Where(&model.PigDeviceSlave{
		DeviceId:      slave.DeviceId,
		ModbusAddress: slaveId,
	}).First(&slave).Error
	if err == nil {
		db.RDB.HSet(db.RDB.GetSlaveKey(iccid), strconv.Itoa(slaveId), slave)
		return getProperty(slave.PropertyId)
	}

	return nil, err
}
func getProperty(propertyId int) (*model.PigProperty, error) {
	var property *model.PigProperty
	err := db.RDB.HGet(db.RDB.GetPropertyKey(), strconv.Itoa(propertyId), &property)
	if err == nil {
		return property, err
	}
	var propertys []*model.PigProperty
	err = db.DB.Find(&propertys).Error
	if err != nil {
		return property, err
	}
	for _, pigProperty := range propertys {
		db.RDB.HSet(db.RDB.GetPropertyKey(), strconv.Itoa(pigProperty.Id), pigProperty)
	}
	err = db.RDB.HGet(db.RDB.GetPropertyKey(), strconv.Itoa(propertyId), &property)
	if err == nil {
		return property, err
	}
	return property, nil
}
func sendAwait5Second(iccid string, slaveId int) bool {
	key := db.RDB.GetAwaitSendKey(iccid, strconv.Itoa(slaveId))
	if !db.RDB.Has(key) {
		db.RDB.Set(key, 1, 5*time.Minute)
		return true
	} else {
		return false
	}
}
func clearSendAwait(iccid string, slaveId int) {
	key := db.RDB.GetAwaitSendKey(iccid, strconv.Itoa(slaveId))
	db.RDB.Del(key)
}
