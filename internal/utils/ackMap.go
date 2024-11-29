package utils

import "sync"

func GetPacketStatus(packetID byte, packetMap map[byte]bool, pMutex *sync.Mutex) (bool, bool) {
	pMutex.Lock()
	waiting, exists := packetMap[packetID]
	pMutex.Unlock()
	return waiting, exists
}

func PacketIsWaiting(packetID byte, packetMap map[byte]bool, pMutex *sync.Mutex, isWaiting bool) {
	pMutex.Lock()
	if _, exists := packetMap[packetID]; !exists {
		pMutex.Unlock()
		return
	}
	packetMap[packetID] = isWaiting
	pMutex.Unlock()
}
