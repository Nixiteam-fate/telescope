package model

import (
	ces_utils "github.com/huaweicloud/telescope/agent/core/ces/utils"
	"github.com/huaweicloud/telescope/agent/core/utils"
)

// GBConversion the multiple of GB --> Byte
const GBConversion = 1024 * 1024 * 1024

// Metric the type for metric data
type Metric struct {
	MetricName   string  `json:"metric_name"`
	MetricValue  float64 `json:"metric_value"`
	ExtraDimension *DimensionType  `json:"extra_dimen,omitempty"`
}

// InputMetric the type for input metric
type InputMetric struct {
	CollectTime int64    `json:"collect_time"`
	Data        []Metric `json:"data"`
}

// InputMetricSlice the type for input metric sclice
type InputMetricSlice []*InputMetric

// DimensionType the type for dimension
type DimensionType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CesMeticExtraInfo struct {
	OriginMetricName string `json:"origin_metric_name"`
	MetricPrefix     string `json:"metric_prefix,omitempty"`
}

// MetricType the type for metric
type MetricType struct {
	Namespace       string             `json:"namespace"`
	Dimensions      []DimensionType    `json:"dimensions"`
	MetricName      string             `json:"metric_name"`
	MetricExtraInfo *CesMeticExtraInfo `json:"extra_info,omitempty"`
}

// CesMetricData the type for post metric data
type CesMetricData struct {
	Metric      MetricType `json:"metric"`
	TTL         int        `json:"ttl"`
	CollectTime int64      `json:"collect_time"`
	Value       float64    `json:"value"`
	Unit        string     `json:"unit"`
}

// CesMetricDataArr the type for metric data array
type CesMetricDataArr []CesMetricData

var metricUnitMap = map[string]string{
	"cpu_usage":					 "%",
	"cpu_usage_user":               "%",
	"cpu_usage_system":             "%",
	"cpu_usage_idle":               "%",
	"cpu_usage_other":              "%",
	"cpu_usage_nice":               "%",
	"cpu_usage_iowait":             "%",
	"cpu_usage_irq":                "%",
	"cpu_usage_softirq":            "%",
	"cpu_usage_steal":              "%",
	"cpu_usage_guest":              "%",
	"cpu_usage_guest_nice":         "%",
	"mem_total":                    "GB",
	"mem_available":                "GB",
	"mem_used":                     "GB",
	"mem_free":                     "GB",
	"mem_usedPercent":              "%",
	"mem_buffers":                  "GB",
	"mem_cached":                   "GB",
	"net_bitSent":                  "bits/s",
	"net_bitRecv":                  "bits/s",
	"net_packetSent":               "Counts/s",
	"net_packetRecv":               "Counts/s",
	"net_errin":                    "%",
	"net_errout":                   "%",
	"net_dropin":                   "%",
	"net_dropout":                  "%",
	"net_fifoin":                   "Bytes",
	"net_fifoout":                  "Bytes",
	"disk_total":                   "GB",
	"disk_free":                    "GB",
	"disk_used":                    "GB",
	"disk_usedPercent":             "%",
	"disk_inodesTotal":             "",
	"disk_inodesUsed":              "",
	"disk_inodesFree":              "",
	"disk_inodesUsedPercent":       "%",
	"disk_writeBytes":              "Bytes",
	"disk_readBytes":               "Bytes",
	"disk_iopsInProgress":          "Bytes",
	"disk_agt_read_bytes_rate":     "Byte/s",
	"disk_agt_read_requests_rate":  "Requests/Second",
	"disk_agt_write_bytes_rate":    "Byte/s",
	"disk_agt_write_requests_rate": "Requests/Second",
	"disk_writeTime":               "ms/Count",
	"disk_readTime":                "ms/Count",
	"disk_ioUtils":                 "%",
	"proc_cpu":                     "%",
	"proc_mem":                     "%",
	"proc_file":                    "Count",
	"gpu_performance_state":        "",
	"gpu_usage_gpu":                "%",
	"gpu_usage_mem":                "%",
	"load_average1":				 "Task/CPU",
	"load_average5":				 "Task/CPU",
	"load_average15":				 "Task/CPU",
}

// BuildMetric build metric as input metric
func BuildMetric(collectTime int64, data []Metric) *InputMetric {
	return &InputMetric{
		CollectTime: collectTime,
		Data:        data,
	}
}

// BuildCesMetricData build ces metric data
func BuildCesMetricData(inputMetric *InputMetric, isAggregated bool) CesMetricDataArr {
	var dimension DimensionType
	var metricTTL int
	var cesMetricDataArr CesMetricDataArr

	dimension.Name = ces_utils.DimensionName
	dimension.Value = utils.GetConfig().InstanceId
	dimensions := []DimensionType{}
	dimensions = append(dimensions, dimension)
	collectTime := inputMetric.CollectTime
	namespace := ces_utils.NameSpace

	externalNamespace := utils.GetConfig().ExternalService
	if externalNamespace == ces_utils.ExternalServiceBMS {
		namespace = externalNamespace
	}

	if isAggregated {
		metricTTL = ces_utils.TTLTwoDay
	} else {
		metricTTL = ces_utils.TTLOneHour
	}

	for _, metric := range inputMetric.Data {

		var newMetricData CesMetricData

		newMetricData.Metric.MetricName = metric.MetricName
		if (metric.ExtraDimension != nil) {
			newMetricData.Metric.Dimensions = append(dimensions, *metric.ExtraDimension)
		} else {
			newMetricData.Metric.Dimensions = dimensions
		}

		newMetricData.Metric.Namespace = namespace
		newMetricData.CollectTime = collectTime
		newMetricData.TTL = metricTTL
		newMetricData.Value = utils.Limit2Decimal(metric.MetricValue)
		newMetricData.Unit = getUnitByMetric(metric.MetricName)

		cesMetricDataArr = append(cesMetricDataArr, newMetricData)

	}
	return cesMetricDataArr

}

func getUnitByMetric(metricName string) string {

	return metricUnitMap[metricName]

}
