package collectors

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/shirou/gopsutil/disk"
)

// DiskCollector is the collector type for disk metric
type DiskCollector struct {
}

const MOUNT_POINT  = "mount_point"
// Collect implement the disk Collector
func (d *DiskCollector) Collect(collectTime int64) *model.InputMetric {

	var result model.InputMetric

	fieldsG := []model.Metric{}

	diskPartitions, _ := disk.Partitions(false)
	diskInfo, _ := disk.IOCounters()

	for _, eachDisk := range diskPartitions {

		diskMountpoint := eachDisk.Mountpoint
		diskStats, _ := disk.Usage(diskMountpoint)

		diskName := eachDisk.Device

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_total", MetricValue: float64(diskStats.Total) / model.GBConversion, ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_free", MetricValue: float64(diskStats.Free) / model.GBConversion, ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_used", MetricValue: float64(diskStats.Used) / model.GBConversion, ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_usedPercent", MetricValue: float64(diskStats.UsedPercent), ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})

		if diskInfo[diskName].Name == "" {
			logs.GetCesLogger().Infof("No IO data for the disk : %v", diskName)
			continue
		}

		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_bytes_rate", MetricValue: float64(diskInfo[diskName].ReadBytes), ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_read_requests_rate", MetricValue: float64(diskInfo[diskName].ReadCount), ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_bytes_rate", MetricValue: float64(diskInfo[diskName].WriteBytes), ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
		fieldsG = append(fieldsG, model.Metric{MetricName: "disk_agt_write_requests_rate", MetricValue: float64(diskInfo[diskName].WriteCount), ExtraDimension:&model.DimensionType{Name:MOUNT_POINT, Value:diskMountpoint}})
	}

	result.Data = fieldsG
	result.CollectTime = collectTime

	return &result
}
