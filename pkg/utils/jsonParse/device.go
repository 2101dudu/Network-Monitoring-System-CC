package jsonParse

type Device struct {
	DeviceID            uint8               `json:"device_id"` // 0 - 255
	DeviceMetrics       DeviceMetrics       `json:"device_metrics"`
	LinkMetrics         LinkMetrics         `json:"link_metrics"`
	AlertFlowConditions AlertFlowConditions `json:"alertflow_conditions"`
}

func validateDevice(device Device, numberOfDevices int) bool {
    if (!validateDeviceMetrics(device.DeviceMetrics) || !validateLinkMetrics(device.LinkMetrics, numberOfDevices) || !validateAlertFlowConditions(device.AlertFlowConditions)) {
        return false
    }
    return true
}
