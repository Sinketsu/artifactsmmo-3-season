package generic

import ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"

var (
	characterLevel = ycmonitoringgo.NewIGauge("level", ycmonitoringgo.DefaultRegistry, "_name")
	skillLevel     = ycmonitoringgo.NewIGauge("skill_level", ycmonitoringgo.DefaultRegistry, "_name", "skill")
	goldCount      = ycmonitoringgo.NewIGauge("gold_count", ycmonitoringgo.DefaultRegistry, "_name")

	apiRequestCount = ycmonitoringgo.NewCounter("api_request_count", ycmonitoringgo.DefaultRegistry)

	itemCount = ycmonitoringgo.NewIGauge("item_count", ycmonitoringgo.DefaultRegistry, "_name", "item")
)
