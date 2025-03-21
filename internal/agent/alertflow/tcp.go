package alertflow

import (
	alertTcp "nms/internal/packet/alert"
	utils "nms/internal/utils"
)

func ConnectTCPAndSendAlert(serverTCPPort string, alert alertTcp.Alert) {

	conn := utils.ResolveTCPAddrAndDial(utils.SERVERIP, serverTCPPort)
	defer conn.Close()

	alertTcp.EncodeAndSendAlert(conn, alert)
}
