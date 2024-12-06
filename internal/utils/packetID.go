package utils

import (
	"sync"
)

func ReadAndIncrementPacketID(packetID *uint16, packetMutex *sync.Mutex, increment bool) uint16 {
	packetMutex.Lock()
	defer packetMutex.Unlock()
	id := *packetID
	if increment {
		*packetID++
	}
	return id
}
