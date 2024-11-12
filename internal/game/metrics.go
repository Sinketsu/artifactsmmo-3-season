package game

import (
	ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"
)

var (
	// TODO fix Rate type
	apiRequestCount = ycmonitoringgo.NewCounter("game_api_request_count", ycmonitoringgo.DefaultRegistry, "_service")

	achievmentsStatus = ycmonitoringgo.NewDGauge("achievments", ycmonitoringgo.DefaultRegistry, "_name")
	bankGoldCount     = ycmonitoringgo.NewIGauge("bank_gold_count", ycmonitoringgo.DefaultRegistry)
	bankItemCount     = ycmonitoringgo.NewIGauge("bank_item_count", ycmonitoringgo.DefaultRegistry, "item")
)
