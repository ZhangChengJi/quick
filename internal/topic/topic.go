package topic

const (

	//属性配置上报 QoS1
	Property_config_post_topic = "sys/+/thing/property/config/post"

	//属性配置主动获取
	Property_config_get_topic = "sys/+/thing/property/config/get"

	//属性上报topic
	Property_post_topic = "sys/+/thing/property/post"

	//事件上报topic QoS1
	Event_post_topic = "sys/+/thing/event/post"

	Device_connect = "sys/+/thing/connect"

	Device_disconnect = "sys/+/thing/disconnect"

	K_device_metadata_chanl = "k_device_metadata_chan"

	K_device_data_chanl = "k_device_data_chan"

	K_device_event_chanl = "k_device_event_chan"

	K_device_line = "k_device_line"

	Device_event  = "device/%s/%s/event/post"
	Device_notify = "device/%s/notify/post"
	Device_line   = "device/%s/%s/line/post"
	Device_last   = "device/%s/property/post"

	OpenApi_data = "api/%s/%s/thing/post"
	OpenApi_line = "api/%s/%s/line/post"
)
