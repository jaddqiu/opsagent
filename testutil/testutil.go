package testutil

import (
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/metric"
)

var localhost = "localhost"

// GetLocalHost returns the DOCKER_HOST environment variable, parsing
// out any scheme or ports so that only the IP address is returned.
func GetLocalHost() string {
	if dockerHostVar := os.Getenv("DOCKER_HOST"); dockerHostVar != "" {
		u, err := url.Parse(dockerHostVar)
		if err != nil {
			return dockerHostVar
		}

		// split out the ip addr from the port
		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			return dockerHostVar
		}

		return host
	}
	return localhost
}

// MockMetrics returns a mock []opsagent.Metric object for using in unit tests
// of opsagent output sinks.
func MockMetrics() []opsagent.Metric {
	metrics := make([]opsagent.Metric, 0)
	// Create a new point batch
	metrics = append(metrics, TestMetric(1.0))
	return metrics
}

// TestMetric Returns a simple test point:
//     measurement -> "test1" or name
//     tags -> "tag1":"value1"
//     value -> value
//     time -> time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
func TestMetric(value interface{}, name ...string) opsagent.Metric {
	if value == nil {
		panic("Cannot use a nil value")
	}
	measurement := "test1"
	if len(name) > 0 {
		measurement = name[0]
	}
	tags := map[string]string{"tag1": "value1"}
	pt, _ := metric.New(
		measurement,
		tags,
		map[string]interface{}{"value": value},
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	)
	return pt
}
