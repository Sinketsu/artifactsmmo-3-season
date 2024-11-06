package game

import (
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

var (
	mapRequestRate = ycmonitoringgo.NewRate("map_request_count", ycmonitoringgo.DefaultRegistry)
)
