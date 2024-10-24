package fluxnode

const (
	topicRequestConfig = "flux.service.%s.request-config"
	topicOnReady       = "flux.service.%s.ready"
	topicOnStart       = "flux.service.%s.start"
	topicOnStop        = "flux.service.%s.stop"
	topicOnRestart     = "flux.service.%s.restart"
	topicOnError       = "flux.service.%s.error"
	topicOnTick        = "flux.service.%s.tick"
)
