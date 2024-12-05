package task

import (
	"bytes"
	"log"
	utils "nms/internal/utils"
)

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

type DeviceMetricsBuilder struct {
	DeviceMetrics DeviceMetrics
}

func NewDeviceMetricsBuilder() *DeviceMetricsBuilder {
	return &DeviceMetricsBuilder{
		DeviceMetrics: DeviceMetrics{
			CpuUsage:       false,
			RamUsage:       false,
			InterfaceStats: []string{},
		},
	}
}

func (b *DeviceMetricsBuilder) SetCpuUsage(usage bool) *DeviceMetricsBuilder {
	b.DeviceMetrics.CpuUsage = usage
	return b
}

func (b *DeviceMetricsBuilder) SetRamUsage(usage bool) *DeviceMetricsBuilder {
	b.DeviceMetrics.RamUsage = usage
	return b
}

func (b *DeviceMetricsBuilder) SetInterfaceStats(stat []string) *DeviceMetricsBuilder {
	b.DeviceMetrics.InterfaceStats = stat
	return b
}

func (b *DeviceMetricsBuilder) Build() DeviceMetrics {
	return b.DeviceMetrics
}

func EncodeDeviceMetrics(metrics DeviceMetrics) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(utils.BoolToByte(metrics.CpuUsage))
	buf.WriteByte(utils.BoolToByte(metrics.RamUsage))
	buf.WriteByte(byte(len(metrics.InterfaceStats)))
	for _, stat := range metrics.InterfaceStats {
		statBytes := []byte(stat)
		buf.WriteByte(byte(len(statBytes)))
		buf.Write(statBytes)
	}

	packet := buf.Bytes()

	if len(packet) > utils.BUFFERSIZE {
		log.Fatalln(utils.Red+"[ERROR 205] Packet size too large", utils.Reset)
	}

	return packet
}

func DecodeDeviceMetrics(data []byte) (DeviceMetrics, error) {
	buf := bytes.NewReader(data)
	var metrics DeviceMetrics

	cpuUsage, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	ramUsage, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	metrics.CpuUsage = cpuUsage == 1
	metrics.RamUsage = ramUsage == 1

	numStats, err := buf.ReadByte()
	if err != nil {
		return metrics, err
	}
	for i := 0; i < int(numStats); i++ {
		statLen, err := buf.ReadByte()
		if err != nil {
			return metrics, err
		}
		statBytes := make([]byte, statLen)
		if _, err := buf.Read(statBytes); err != nil {
			return metrics, err
		}
		metrics.InterfaceStats = append(metrics.InterfaceStats, string(statBytes))
	}
	return metrics, nil
}
