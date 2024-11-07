package game

import (
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

var (
	apiRequestCount = ycmonitoringgo.NewRate("game_api_request_count", ycmonitoringgo.DefaultRegistry, "service")
)
