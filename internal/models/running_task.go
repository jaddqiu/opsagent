package models

import (
	"time"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/selfstat"
)

//var GlobalMetricsGathered = selfstat.Register("agent", "metrics_gathered", map[string]string{})

type RunningTask struct {
	Task   opsagent.Task
	Config *TaskConfig

	defaultTags map[string]string

	MetricsGathered selfstat.Stat
	GatherTime      selfstat.Stat
}

func NewRunningTask(task opsagent.Task, config *TaskConfig) *RunningTask {
	return &RunningTask{
		Task:   task,
		Config: config,
		MetricsGathered: selfstat.Register(
			"gather",
			"metrics_gathered",
			map[string]string{"task": config.Name},
		),
		GatherTime: selfstat.RegisterTiming(
			"gather",
			"gather_time_ns",
			map[string]string{"task": config.Name},
		),
	}
}

// TaskConfig is the common config for all tasks.
type TaskConfig struct {
	Name     string
	Interval time.Duration

	CronSpec string
	Schedule string
	Filter   Filter
}

func (r *RunningTask) Name() string {
	return "tasks." + r.Config.Name
}

func (r *RunningTask) metricFiltered(metric opsagent.Metric) {
	metric.Drop()
}

/*
func (r *RunningTask) MakeMetric(metric opsagent.Metric) opsagent.Metric {
	if ok := r.Config.Filter.Select(metric); !ok {
		r.metricFiltered(metric)
		return nil
	}

	m := makemetric(
		metric,
		r.Config.MeasurementPrefix,
		r.Config.MeasurementSuffix,
		r.Config.Tags,
		r.defaultTags)

	r.Config.Filter.Modify(metric)
	if len(metric.FieldList()) == 0 {
		r.metricFiltered(metric)
		return nil
	}

	r.MetricsGathered.Incr(1)
	GlobalMetricsGathered.Incr(1)
	return m
}
*/

func (r *RunningTask) Execute() error {
	err := r.Task.Execute()
	return err
}
