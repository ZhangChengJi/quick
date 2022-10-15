package model

type PigDeviceSlave struct {
	Id            int64  //
	DeviceId      string // 设备ID
	SlaveAlias    string // 从机设备别名
	SlaveName     string // 从机设备名称
	ModbusAddress int    // modbus从站地址
	PropertyId    int    // 关联设备属性
	SlaveDesc     string // 从机设备描述
	SlaveStatus   int    // 从机设备开关 0:关闭 1:开启
	LineStatus    int    // 从机设备状态 0:离线 1:在线
}

func (PigDeviceSlave) TableName() string {
	return "pig_device_slave"
}
