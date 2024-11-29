package metrics

import (
	"errors"
	"nms/internal/utils"
)

type Metrics struct {
	PacketID byte
	AgentID  byte
	Metrics  string
}

type MetricsBuilder struct {
	Metrics Metrics
}

func NewMetricsBuilder() *MetricsBuilder {
	return &MetricsBuilder{
		Metrics: Metrics{
			PacketID: 0,
			AgentID:  0,
			Metrics:  "",
		},
	}
}

func (m *MetricsBuilder) SetPacketID(packetID byte) *MetricsBuilder {
	m.Metrics.PacketID = packetID
	return m
}

func (m *MetricsBuilder) SetAgentID(agentID byte) *MetricsBuilder {
	m.Metrics.AgentID = agentID
	return m
}

func (m *MetricsBuilder) SetMetrics(metrics string) *MetricsBuilder {
	m.Metrics.Metrics = metrics
	return m
}

func (m *MetricsBuilder) Build() Metrics {
	return m.Metrics
}

func DecodeMetrics(packet []byte) (Metrics, error) {
	if len(packet) < 2 {
		return Metrics{}, errors.New("invalid packet length")
	}

	return Metrics{
		PacketID: packet[0],
		AgentID:  packet[1],
		Metrics:  string(packet[2:]),
	}, nil
}

func EncodeMetrics(metrics Metrics) []byte {
	return append([]byte{byte(utils.METRICSGATHERING), metrics.PacketID, metrics.AgentID}, []byte(metrics.Metrics)...)
}
