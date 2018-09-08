package aggregate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// AvgValue is the average result type for Aggregate
type AvgValue struct {
}

// Aggregate implement the average aggregator
func (averageValue *AvgValue) Aggregate(input model.InputMetricSlice) *model.InputMetric {

	if input == nil || len(input) == 0 {
		logs.GetCesLogger().Error("Input slice is nil or empty")
		return nil
	}

	dataCount := len(input)
	avgMetric := *input[0]

	//aggregate collectTime Round to Minute
	avgMetric.CollectTime = time.Unix(avgMetric.CollectTime/1000, 0).Truncate(time.Minute).Unix() * 1000
	metricNameKeyMap := GenerateMetricNameKeyMap(&avgMetric.Data)

	metricCount := len((*input[0]).Data)

	sum := make(map[string]float64, metricCount)
	prefix := ""
	for _, metricData := range input {

		for _, metric := range metricData.Data {
			if (metric.ExtraDimension != nil){
				prefix = metric.ExtraDimension.Value
			}
			sum[prefix + metric.MetricName] = sum[prefix + metric.MetricName] + metric.MetricValue
		}

	}
	prefix = ""
	for _, metric := range avgMetric.Data {
		if (metric.ExtraDimension != nil){
			prefix = metric.ExtraDimension.Value
		}
		avg := sum[prefix + metric.MetricName]/float64(dataCount)
		metricNameKeyMap[prefix + metric.MetricName].MetricValue, _ = strconv.ParseFloat(fmt.Sprint("%.2f", avg), 64)
	}

	return &avgMetric

}
