package metrics

import (
	"bytes"
	"encoding/binary"
	"errors"
	"nms/internal/utils"
)

type Metrics struct {
	PacketID byte
	AgentID  byte
	TaskID   uint16
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
			TaskID:   0,
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

func (m *MetricsBuilder) SetTaskID(taskID uint16) *MetricsBuilder {
	m.Metrics.TaskID = taskID
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
	buf := bytes.NewReader(packet)
	var metrics Metrics

	if len(packet) < 7 {
		return metrics, errors.New("invalid packet length")
	}

	packetID, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	agentID, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	if err := binary.Read(buf, binary.BigEndian, &metrics.TaskID); err != nil {
		return metrics, err
	}
	metrics.PacketID = packetID
	metrics.AgentID = agentID

	var metricsLen uint16
	if err := binary.Read(buf, binary.BigEndian, &metricsLen); err != nil {
		return metrics, err
	}

	metricsBytes := make([]byte, metricsLen)
	if _, err := buf.Read(metricsBytes); err != nil {
		return metrics, err
	}
	metrics.Metrics = string(metricsBytes)

	return metrics, nil
}

func EncodeMetrics(metrics Metrics) []byte {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(utils.METRICSGATHERING))
	buf.WriteByte(metrics.PacketID)
	buf.WriteByte(metrics.AgentID)
	binary.Write(buf, binary.BigEndian, metrics.TaskID)

	metricsBytes := []byte(metrics.Metrics)
	metricsLen := uint16(len(metricsBytes))
	binary.Write(buf, binary.BigEndian, metricsLen)
	buf.Write(metricsBytes)

	return buf.Bytes()
}
