package all

import (
	_ "github.com/jaddqiu/opsagent/plugins/outputs/amon"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/amqp"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/application_insights"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/azure_monitor"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/cloud_pubsub"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/cloudwatch"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/cratedb"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/datadog"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/discard"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/elasticsearch"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/file"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/graphite"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/graylog"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/http"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/influxdb"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/influxdb_v2"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/instrumental"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/kafka"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/kinesis"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/librato"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/mqtt"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/nats"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/nsq"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/opentsdb"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/prometheus_client"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/riemann"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/riemann_legacy"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/socket_writer"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/stackdriver"
	_ "github.com/jaddqiu/opsagent/plugins/outputs/wavefront"
)
