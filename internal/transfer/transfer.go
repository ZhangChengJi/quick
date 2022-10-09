package transfer

import (
	"database/sql"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/**
Topic列表
物联网平台预定义物模型通信Topic，各物模型功能Topic消息的数据格式，见OneJSON数据格式。
Topic类以正斜线（/）进行分层，区分每个类目。其中，有两个类目为既定类目：${pid}表示产品的产品id；${device-name}表示设备名称；${identifier}表示服务标识符

功能	类别	行为	描述	Topic类	操作权限
物模型通信Topic	属性	设备属性上报	请求	$sys/{pid}/{device-name}/thing/property/post	发布
响应	$sys/{pid}/{device-name}/thing/property/post/reply	订阅
设备属性设置（同步）	请求	$sys/{pid}/{device-name}/thing/property/set	订阅
响应	$sys/{pid}/{device-name}/thing/property/set_reply	发布
设备获取属性期望值	请求	$sys/{pid}/{device-name}/thing/property/desired/get	发布
响应	$sys/{pid}/{device-name}/thing/property/desired/get/reply	订阅
清除属性期望值	请求	$sys/{pid}/{device-name}/thing/property/desired/delete	发布
响应	$sys/{pid}/{device-name}/thing/property/desired/delete/reply	订阅
设备属性获取	请求	$sys/{pid}/{device-name}/thing/property/get	订阅
响应	$sys/{pid}/{device-name}/thing/property/get_reply	发布
事件	设备事件上报	请求	$sys/{pid}/{device-name}/thing/event/post	发布
响应	$sys/{pid}/{device-name}/thing/event/post/reply	订阅
服务	设备服务调用	请求	$sys/{pid}/{device-name}/thing/service/{identifier}/invoke	订阅
响应	$sys/{pid}/{device-name}/thing/service/{identifier}/invoke_reply	发布
功能	类别	行为	描述	Topic类	操作权限
网关与子设备通信Topic	上下线	子设备上线	请求	$sys/{pid}/{device-name}/thing/sub/login	发布
响应	$sys/{pid}/{device-name}/thing/sub/login/reply	订阅
子设备下线	请求	$sys/{pid}/{device-name}/thing/sub/logout	发布
响应	$sys/{pid}/{device-name}/thing/sub/logout/reply	订阅
属性和事件	批量上报属性和事件（网关上报或代理子设备上报）	请求	$sys/{pid}/{device-name}/thing/pack/post	发布
响应	$sys/{pid}/{device-name}/thing/pack/post/reply	订阅
子设备属性获取	请求	$sys/{pid}/{device-name}/thing/sub/property/get	订阅
响应	$sys/{pid}/{device-name}/thing/sub/property/get_reply	发布
子设备属性设置	请求	$sys/{pid}/{device-name}/thing/sub/property/set	订阅
响应	$sys/{pid}/{device-name}/thing/sub/property/set_reply	发布
历史属性和事件上报（网关上报或代理子设备上报）	请求	$sys/{pid}/{device-name}/thing/history/post	发布
响应	$sys/{pid}/{device-name}/thing/history/post/reply	订阅
服务	子设备服务调用	请求	$sys/{pid}/{device-name}/thing/sub/service/invoke	订阅
响应	$sys/{pid}/{device-name}/thing/sub/service/invoke_reply	发布
拓扑关系	添加子设备	请求	$sys/{pid}/{device-name}/thing/sub/topo/add	发布
响应	$sys/{pid}/{device-name}/thing/sub/topo/add/reply	订阅
删除子设备	请求	$sys/{pid}/{device-name}/thing/sub/topo/delete	发布
响应	$sys/{pid}/{device-name}/thing/sub/topo/delete/reply	订阅
获取拓扑关系	请求	$sys/{pid}/{device-name}/thing/sub/topo/get	发布
响应	$sys/{pid}/{device-name}/thing/sub/topo/get/reply	订阅
网关同步结果响应	$sys/{pid}/{device-name}/thing/sub/topo/get/result	发布
通知网关拓扑关系变化	请求	$sys/{pid}/{device-name}/thing/sub/topo/change	订阅
响应	$sys/{pid}/{device-name}/thing/sub/topo/change_reply	发布
*/
const (
	property_topic = "sys/+/thing/property/post" //设备属性上报
)

type transfer struct {
	mq mqtt.Client
	td *sql.DB
}

func Start() {

}
func (t *transfer) listenMqtt() {
	t.mq.Subscribe(property_topic, 0, t.stream)
}
func (*transfer) stream(client mqtt.Client, msg mqtt.Message) {

}
