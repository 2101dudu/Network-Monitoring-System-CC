package packet

import "sync"



func GetPacketIDStatus(packetID byte, packetMap map[byte]bool, pMutex *sync.Mutex) (bool, bool)  {
    pMutex.Lock()
    waiting, exists := packetMap[packetID]
    pMutex.Unlock()
    return waiting, exists
}

func PacketIDIsWaiting(packetID byte, packetMap map[byte]bool, pMutex *sync.Mutex, isWaiting bool) {
	pMutex.Lock()
	packetMap[packetID] = isWaiting
	pMutex.Unlock()
}
