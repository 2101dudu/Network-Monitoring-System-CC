package nettask

import (
	parse "nms/internal/jsonParse"
)

func getIperfIndices(task parse.Task) (byte, byte) {
	if task.Devices[0].LinkMetrics.IperfParameters.IsServer {
		return 0, 1
	}
	return 1, 0
}
