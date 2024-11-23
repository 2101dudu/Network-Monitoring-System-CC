package task

import (
	"bytes"
	utils "nms/pkg/utils"
)

type DeviceMetrics struct {
	CpuUsage       bool
	RamUsage       bool
	InterfaceStats []string
}

func EncodeDeviceMetrics(metrics DeviceMetrics) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(utils.BoolToByte(metrics.CpuUsage))
	buf.WriteByte(utils.BoolToByte(metrics.RamUsage))
	buf.WriteByte(byte(len(metrics.InterfaceStats)))
	for _, stat := range metrics.InterfaceStats {
		statBytes := []byte(stat)
		buf.WriteByte(byte(len(statBytes)))
		buf.Write(statBytes)
	}
	return buf.Bytes(), nil
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
