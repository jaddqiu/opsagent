package metric

import (
	"testing"
	"time"

	"github.com/jaddqiu/opsagent"
	"github.com/stretchr/testify/require"
)

func mustMetric(
	name string,
	tags map[string]string,
	fields map[string]interface{},
	tm time.Time,
	tp ...opsagent.ValueType,
) opsagent.Metric {
	m, err := New(name, tags, fields, tm, tp...)
	if err != nil {
		panic("mustMetric")
	}
	return m
}

type deliveries struct {
	Info map[opsagent.TrackingID]opsagent.DeliveryInfo
}

func (d *deliveries) onDelivery(info opsagent.DeliveryInfo) {
	d.Info[info.ID()] = info
}

func TestTracking(t *testing.T) {
	tests := []struct {
		name      string
		metric    opsagent.Metric
		actions   func(metric opsagent.Metric)
		delivered bool
	}{
		{
			name: "accept",
			metric: mustMetric(
				"cpu",
				map[string]string{},
				map[string]interface{}{
					"value": 42,
				},
				time.Unix(0, 0),
			),
			actions: func(m opsagent.Metric) {
				m.Accept()
			},
			delivered: true,
		},
		{
			name: "reject",
			metric: mustMetric(
				"cpu",
				map[string]string{},
				map[string]interface{}{
					"value": 42,
				},
				time.Unix(0, 0),
			),
			actions: func(m opsagent.Metric) {
				m.Reject()
			},
			delivered: false,
		},
		{
			name: "accept copy",
			metric: mustMetric(
				"cpu",
				map[string]string{},
				map[string]interface{}{
					"value": 42,
				},
				time.Unix(0, 0),
			),
			actions: func(m opsagent.Metric) {
				m2 := m.Copy()
				m.Accept()
				m2.Accept()
			},
			delivered: true,
		},
		{
			name: "copy with accept and done",
			metric: mustMetric(
				"cpu",
				map[string]string{},
				map[string]interface{}{
					"value": 42,
				},
				time.Unix(0, 0),
			),
			actions: func(m opsagent.Metric) {
				m2 := m.Copy()
				m.Accept()
				m2.Drop()
			},
			delivered: true,
		},
		{
			name: "copy with mixed delivery",
			metric: mustMetric(
				"cpu",
				map[string]string{},
				map[string]interface{}{
					"value": 42,
				},
				time.Unix(0, 0),
			),
			actions: func(m opsagent.Metric) {
				m2 := m.Copy()
				m.Accept()
				m2.Reject()
			},
			delivered: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &deliveries{
				Info: make(map[opsagent.TrackingID]opsagent.DeliveryInfo),
			}
			metric, id := WithTracking(tt.metric, d.onDelivery)
			tt.actions(metric)

			info := d.Info[id]
			require.Equal(t, tt.delivered, info.Delivered())
		})
	}
}

func TestGroupTracking(t *testing.T) {
	tests := []struct {
		name      string
		metrics   []opsagent.Metric
		actions   func(metrics []opsagent.Metric)
		delivered bool
	}{
		{
			name: "accept",
			metrics: []opsagent.Metric{
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
			},
			actions: func(metrics []opsagent.Metric) {
				metrics[0].Accept()
				metrics[1].Accept()
			},
			delivered: true,
		},
		{
			name: "reject",
			metrics: []opsagent.Metric{
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
			},
			actions: func(metrics []opsagent.Metric) {
				metrics[0].Reject()
				metrics[1].Reject()
			},
			delivered: false,
		},
		{
			name: "remove",
			metrics: []opsagent.Metric{
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
			},
			actions: func(metrics []opsagent.Metric) {
				metrics[0].Drop()
				metrics[1].Drop()
			},
			delivered: true,
		},
		{
			name: "mixed",
			metrics: []opsagent.Metric{
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
				mustMetric(
					"cpu",
					map[string]string{},
					map[string]interface{}{
						"value": 42,
					},
					time.Unix(0, 0),
				),
			},
			actions: func(metrics []opsagent.Metric) {
				metrics[0].Accept()
				metrics[1].Reject()
			},
			delivered: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &deliveries{
				Info: make(map[opsagent.TrackingID]opsagent.DeliveryInfo),
			}
			metrics, id := WithGroupTracking(tt.metrics, d.onDelivery)
			tt.actions(metrics)

			info := d.Info[id]
			require.Equal(t, tt.delivered, info.Delivered())
		})
	}
}
