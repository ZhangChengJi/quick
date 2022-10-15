package property

import (
	"quick/manager/model"
	"quick/pkg/log"
	"strconv"
	"time"
)

type service struct {
}

func queryDevice(iccid string) (*model.PigDevice, error) {
	var pigDevice *model.PigDevice
	if err := db.RDB.Get(iccid, &pigDevice); err == nil {
		return pigDevice, nil

	}
	err := db.DB.Where(&model.PigDevice{Id: iccid}).First(&pigDevice).Error
	if err == nil {
		marshal, err := json.Marshal(&pigDevice)
		if err != nil {
			return nil, err
		}
		db.RDB.Set(iccid, string(marshal), -1)
		return pigDevice, nil
	}

	return pigDevice, err
}
func createOrUpdateSlave(slave *model.PigDeviceSlave) {
	err := db.DB.Where(&model.PigDeviceSlave{
		DeviceId:      slave.DeviceId,
		ModbusAddress: slave.ModbusAddress,
	}).Assign(
		&model.PigDeviceSlave{
			DeviceId:      slave.DeviceId,
			ModbusAddress: slave.ModbusAddress,
			PropertyId:    slave.PropertyId,
			SlaveName:     slave.SlaveName,
		},
	).FirstOrCreate(&slave).Error
	if err != nil {
		log.Sugar.Errorf("探测器更新失败%s", err)
		return
	}
	db.RDB.HSet(db.RDB.GetSlaveKey(slave.DeviceId), strconv.Itoa(slave.ModbusAddress), &slave)

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
