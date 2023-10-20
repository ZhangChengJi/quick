package model

type PigDeviceGroup struct {
	DeviceId string `gorm:"column:device_id" json:"deviceId"` // 产品id
	GroupId  int    `gorm:"column:group_id" json:"groupId"`   // 分组id

}

func (PigDeviceGroup) TableName() string {
	return "pig_device_group"
}
