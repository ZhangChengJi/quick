package model

type PigDevice struct {
	Id               string `gorm:"column:id" json:"id"`                              // 设备编号
	ProductId        int    `gorm:"column:product_id" json:"productId"`               // 产品id
	DeptId           int    `gorm:"column:dept_id" json:"deptId"`                     // 部门id
	DeviceName       string `gorm:"column:device_name" json:"deviceName"`             // 设备名称
	NetworkFlag      int    `gorm:"column:network_flag" json:"networkFlag"`           // 联网方式 1:4G-DTU 2:NB-IOT
	InstructFlag     int    `gorm:"column:instruct_flag" json:"instructFlag"`         // 指令下发方式 1:单条下发 2:多条下发
	SimCode          string `gorm:"column:sim_code" json:"simCode"`                   // SIM卡
	BindStatus       int    `gorm:"column:bind_status" json:"bindStatus"`             // 绑定状态 0:未绑定 1:已绑定
	LineStatus       int    `gorm:"column:line_status" json:"lineStatus"`             // 设备状态: 0:离线 1:在线
	DeviceAddress    string `gorm:"column:device_address" json:"deviceAddress"`       // 设备地址
	DeviceCoordinate string `gorm:"column:device_coordinate" json:"deviceCoordinate"` // 设备坐标信息
	GroupId          int    `gorm:"column:group_id" json:"groupId"`                   // 绑定分组id

}

func (PigDevice) TableName() string {
	return "pig_device"
}
