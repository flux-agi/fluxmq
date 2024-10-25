package fluxnode

const (
	topicRequestConfig = "service/%s/get_config"
	topicOnReady       = "service/%s/set_config"
	topicOnStart       = "service/%s/start"
	topicOnStop        = "service/%s/stop"
	topicOnRestart     = "service/%s/restart"
	topicOnError       = "service/%s/error"
	topicOnTick        = "service/%s/tick"
	topicStatus        = "service/%s/status"         // for push status
	topicStatusRequest = "service/%s/request_status" // for request status
)
