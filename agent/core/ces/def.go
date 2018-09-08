package ces

import (
	"github.com/huaweicloud/telescope/agent/core/ces/model"
)

// common variables (chans and vars)
var (
	// Channels
	chAgResult							chan *model.InputMetric
	chAgRawData                         chan model.InputMetricSlice
)

// Initialize the aggregate data channel
func initchAgRawData() {
	chAgRawData = make(chan model.InputMetricSlice, 100)
}

// Get the data channel
func getchAgRawData() chan model.InputMetricSlice {
	if chAgRawData == nil {
		initchAgRawData()
	}

	return chAgRawData
}

// Initialize the agResult channel
func initchAgResult() {
	chAgResult = make(chan *model.InputMetric, 100)
}

// Get the agResult channel
func getchAgResult() chan *model.InputMetric {
	if chAgResult == nil {
		initchAgResult()
	}

	return chAgResult
}