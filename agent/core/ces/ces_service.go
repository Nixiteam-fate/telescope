package ces

import (
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/service"
)

// Service is one of the services of agent
type Service struct {
}

// Init ces Service config and channel
func (s *Service) Init() {

	config.InitConfig()
	config.InitPluginConfig()
	initchAgRawData()
	initchAgResult()

}

// Start make work goroutines
func (s *Service) Start() {
	go services.StartMetricCollectTask(getchAgRawData())
	go services.StartAggregateTask(getchAgResult(), getchAgRawData())
	go services.SendMetricTask(getchAgResult())
}
