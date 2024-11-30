package utils

import (
	"sync"
)

func ReadAndIncrementPacketID(packetID *byte, packetMutex *sync.Mutex, increment bool) byte {
	packetMutex.Lock()
	defer packetMutex.Unlock()
	id := *packetID
	if increment {
		*packetID++
	}
	return id
}
