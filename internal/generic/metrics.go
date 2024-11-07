package generic

import ycmonitoringgo "github.com/Sinketsu/yc-monitoring-go"

var (
	characterLevel = ycmonitoringgo.NewIGauge("level", ycmonitoringgo.DefaultRegistry, "name")
	skillLevel     = ycmonitoringgo.NewIGauge("skill_level", ycmonitoringgo.DefaultRegistry, "name", "skill")

	apiRequestCount = ycmonitoringgo.NewRate("api_request_count", ycmonitoringgo.DefaultRegistry)
)
