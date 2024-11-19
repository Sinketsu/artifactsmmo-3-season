package game

import (
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

var (
	// TODO fix Rate type
	apiRequestCount = ycmonitoringgo.NewCounter("game_api_request_count", ycmonitoringgo.DefaultRegistry, "_service")

	achievmentsStatus        = ycmonitoringgo.NewDGauge("achievments", ycmonitoringgo.DefaultRegistry, "_name")
	achievmentsPointsCurrent = ycmonitoringgo.NewDGauge("achievments_points_current", ycmonitoringgo.DefaultRegistry)
	achievmentsPointsTotal   = ycmonitoringgo.NewDGauge("achievments_points_total", ycmonitoringgo.DefaultRegistry)
	bankGoldCount            = ycmonitoringgo.NewIGauge("bank_gold_count", ycmonitoringgo.DefaultRegistry)
	bankItemCount            = ycmonitoringgo.NewIGauge("bank_item_count", ycmonitoringgo.DefaultRegistry, "item")
)
