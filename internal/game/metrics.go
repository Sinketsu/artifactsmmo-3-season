package game

import (
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

var (
	apiRequestCount = ycmonitoringgo.NewRate("game_api_request_count", ycmonitoringgo.DefaultRegistry, "service")

	achievmentsStatus = ycmonitoringgo.NewDGauge("achievments", ycmonitoringgo.DefaultRegistry, "name")
	bankGoldCount     = ycmonitoringgo.NewIGauge("bank_gold_count", ycmonitoringgo.DefaultRegistry)
)
