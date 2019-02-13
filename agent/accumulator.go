package agent

import (
	"log"
	"time"

	"github.com/jaddqiu/opsagent"
	"github.com/jaddqiu/opsagent/metric"
	"github.com/jaddqiu/opsagent/selfstat"
)

var (
	NErrors = selfstat.Register("agent", "gather_errors", map[string]string{})
)

type MetricMaker interface {
	Name() string
	MakeMetric(metric opsagent.Metric) opsagent.Metric
}

type accumulator struct {
	maker     MetricMaker
	metrics   chan<- opsagent.Metric
	precision time.Duration
}

func NewAccumulator(
	maker MetricMaker,
	metrics chan<- opsagent.Metric,
) opsagent.Accumulator {
	acc := accumulator{
		maker:     maker,
		metrics:   metrics,
		precision: time.Nanosecond,
	}
	return &acc
}

func (ac *accumulator) AddFields(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, opsagent.Untyped, t...)
}

func (ac *accumulator) AddGauge(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, opsagent.Gauge, t...)
}

func (ac *accumulator) AddCounter(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, opsagent.Counter, t...)
}

func (ac *accumulator) AddSummary(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, opsagent.Summary, t...)
}

func (ac *accumulator) AddHistogram(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t ...time.Time,
) {
	ac.addFields(measurement, tags, fields, opsagent.Histogram, t...)
}

func (ac *accumulator) AddMetric(m opsagent.Metric) {
	m.SetTime(m.Time().Round(ac.precision))
	if m := ac.maker.MakeMetric(m); m != nil {
		ac.metrics <- m
	}
}

func (ac *accumulator) addFields(
	measurement string,
	tags map[string]string,
	fields map[string]interface{},
	tp opsagent.ValueType,
	t ...time.Time,
) {
	m, err := metric.New(measurement, tags, fields, ac.getTime(t), tp)
	if err != nil {
		return
	}
	if m := ac.maker.MakeMetric(m); m != nil {
		ac.metrics <- m
	}
}

// AddError passes a runtime error to the accumulator.
// The error will be tagged with the plugin name and written to the log.
func (ac *accumulator) AddError(err error) {
	if err == nil {
		return
	}
	NErrors.Incr(1)
	log.Printf("E! [%s]: Error in plugin: %v", ac.maker.Name(), err)
}

func (ac *accumulator) SetPrecision(precision, interval time.Duration) {
	if precision > 0 {
		ac.precision = precision
		return
	}
	switch {
	case interval >= time.Second:
		ac.precision = time.Second
	case interval >= time.Millisecond:
		ac.precision = time.Millisecond
	case interval >= time.Microsecond:
		ac.precision = time.Microsecond
	default:
		ac.precision = time.Nanosecond
	}
}

func (ac *accumulator) getTime(t []time.Time) time.Time {
	var timestamp time.Time
	if len(t) > 0 {
		timestamp = t[0]
	} else {
		timestamp = time.Now()
	}
	return timestamp.Round(ac.precision)
}

func (ac *accumulator) WithTracking(maxTracked int) opsagent.TrackingAccumulator {
	return &trackingAccumulator{
		Accumulator: ac,
		delivered:   make(chan opsagent.DeliveryInfo, maxTracked),
	}
}

type trackingAccumulator struct {
	opsagent.Accumulator
	delivered chan opsagent.DeliveryInfo
}

func (a *trackingAccumulator) AddTrackingMetric(m opsagent.Metric) opsagent.TrackingID {
	dm, id := metric.WithTracking(m, a.onDelivery)
	a.AddMetric(dm)
	return id
}

func (a *trackingAccumulator) AddTrackingMetricGroup(group []opsagent.Metric) opsagent.TrackingID {
	db, id := metric.WithGroupTracking(group, a.onDelivery)
	for _, m := range db {
		a.AddMetric(m)
	}
	return id
}

func (a *trackingAccumulator) Delivered() <-chan opsagent.DeliveryInfo {
	return a.delivered
}

func (a *trackingAccumulator) onDelivery(info opsagent.DeliveryInfo) {
	select {
	case a.delivered <- info:
	default:
		// This is a programming error in the input.  More items were sent for
		// tracking than space requested.
		panic("channel is full")
	}
}
