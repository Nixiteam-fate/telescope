package services

import (
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/aggregate"
	"github.com/huaweicloud/telescope/agent/core/ces/collectors"
	"github.com/huaweicloud/telescope/agent/core/ces/config"
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/ces/report"
	cesUtils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

func getPlugins()[]*model.PluginCommand{
	if !config.GetConfig().Enable || !config.GetConfig().EnablePlugin{
		return nil
	}

	if config.GetPluginConfig() == nil{
		return nil
	}

	plugins := config.GetPluginConfig().Plugins

	if len(plugins) > cesUtils.MaxPluginNum{
		plugins = plugins[:cesUtils.MaxPluginNum]
	}
	return plugins
}

// StartMetricCollectTask cron job for metric collect
func StartMetricCollectTask(agData chan model.InputMetricSlice) {

	var collectorList []collectors.CollectorInterface

	// simultaneously modify collectorNum in StartAggregateTask when modify the length of collectorList
	collectorList = append(collectorList, &collectors.CPUCollector{})
	collectorList = append(collectorList, &collectors.MemCollector{})
	collectorList = append(collectorList, &collectors.DiskCollector{})
	collectorList = append(collectorList, &collectors.NetCollector{})
	collectorList = append(collectorList, &collectors.LoadCollector{})

	metricSliceArr := make([]model.InputMetricSlice, len(collectorList))

	counter := 0

	time.Sleep(time.Duration(5) * time.Second)

	cronTime := utils.DETAIL_DATA_CRON_JOB_TIME_SECOND//set default
	ticker := time.NewTicker(time.Duration(cronTime) * time.Second)

	for _ = range ticker.C{
		if config.GetConfig().Enable {
			collectTime := time.Now().Unix() * 1000
			for i, collector := range collectorList {

				tmp := collector.Collect(collectTime)

				if tmp != nil {
					metricSliceArr[i] = append(metricSliceArr[i], tmp)
				}
			}
			counter++

			if counter == 6 {
				for i, eachMetricSlice := range metricSliceArr {
					agData <- eachMetricSlice
					metricSliceArr[i] = metricSliceArr[i][:0]
				}
				metricSliceArr = metricSliceArr[:len(collectorList)]
				counter = 0
			}
		}
	}
}

// StartAggregateTask task for aggregate metric in 1 minute
func StartAggregateTask(agRes chan *model.InputMetric, agData chan model.InputMetricSlice) {

	var aggregatorList []aggregate.AggregatorInterface

	aggregatorList = append(aggregatorList, &aggregate.AvgValue{})
	// don't open, first we only support average
	// aggregatorList = append(aggregatorList, &aggregate.MaxValue{})
	// aggregatorList = append(aggregatorList, &aggregate.MinValue{})

	plugins := getPlugins()
	allMetric := new(model.InputMetric)
	tmpData := []model.Metric{}
	count := 0
	collectorNum := 5

	for {
		tmp := <-agData
		for _, aggregator := range aggregatorList {

			eachRes := aggregator.Aggregate(tmp)

			if eachRes != nil {
				for _, value := range eachRes.Data {
					tmpData = append(tmpData, value)
				}

				if allMetric.CollectTime == 0 {
					allMetric.CollectTime = eachRes.CollectTime
				}
			}

		}
		// count length is the collectorNum-1, now the num of collector is 5, and the enabled processes should be considered
		if count >= collectorNum - 1 {
			if plugins != nil {
				for _, plugins := range plugins {
					tmp := plugins.PluginCmd()
					if (tmp != nil) {
						for _, value := range tmp.Data {
							tmpData = append(tmpData, value)
						}
					}
				}
			}
			allMetric.Data = tmpData
			agRes <- allMetric
			tmpData = []model.Metric{}
			count = 0
			allMetric = new(model.InputMetric)
		} else {
			count ++
			continue
		}

	}

}

// BuildURL build URL string by URI
func BuildURL(destURI string) string {
	var url string
	url = config.GetConfig().Endpoint + "/" + utils.API_CES_VERSION + "/" + utils.GetConfig().ProjectId + destURI
	return url
}

// SendMetricTask task for post metric data
func SendMetricTask(agRes chan *model.InputMetric) {
	for {
		metricDataAggregate := <-agRes
		time.Sleep(5 * time.Second)
		go report.SendMetricData(BuildURL(cesUtils.PostAggregatedMetricDataURI), metricDataAggregate, true)
	}
}
