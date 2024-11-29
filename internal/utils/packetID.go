package utils

import (
	"sync"
)

func ReadAndIncrementPacketID(data *byte, packetMutex *sync.Mutex, increment bool) byte {
	packetMutex.Lock()
	defer packetMutex.Unlock()
	id := *data
	if increment {
		*data++
	}
	return id
}
