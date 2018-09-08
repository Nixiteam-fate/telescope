package aggregate

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
)

// MaxValue is the max result type for Aggregate
type MaxValue struct {
}

// Aggregate implement the max aggregator
func (maxValue *MaxValue) Aggregate(input model.InputMetricSlice) *model.InputMetric {

	if input == nil || len(input) == 0 {
		logs.GetCesLogger().Error("Input slice is nil or empty")
		return nil
	}
	maxMetric := *input[0]

	metricNameKeyMap := GenerateMetricNameKeyMap(&maxMetric.Data)
	prefix := ""
	for _, metricData := range input {

		for _, metric := range metricData.Data {

			if (metric.ExtraDimension != nil){
				prefix = metric.ExtraDimension.Value
			}
			if metric.MetricValue > metricNameKeyMap[prefix + metric.MetricName].MetricValue{
				metricNameKeyMap[prefix + metric.MetricName].MetricValue = metric.MetricValue
			}
		}

	}

	return &maxMetric

}

func GenerateMetricNameKeyMap(metrics *[]model.Metric) map[string]*model.Metric {

	metricNameKeyMap := make(map[string]*model.Metric, 0)
	prefix := ""
	for index, _ := range *metrics {
		if ((*metrics)[index].ExtraDimension != nil){
			prefix = (*metrics)[index].ExtraDimension.Value
		}
		metricNameKeyMap[prefix + (*metrics)[index].MetricName] = &(*metrics)[index]
	}

	return metricNameKeyMap
}
