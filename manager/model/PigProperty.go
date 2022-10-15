package model

type PigProperty struct {
	Id                     int    //
	GroupId                int64  // 分类id
	AlarmRule              string // 告警条件
	PropertyName           string // 属性名称
	PropertyIdentification string // 属性标识
	PropertyUnit           string // 单位
	PropertyImg            string // 属性图标
	PropertyDesc           string // 属性描述
}

func (PigProperty) TableName() string {
	return "pig_property"
}
