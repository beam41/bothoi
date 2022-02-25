package gateway

import "sync"

var sequenceNumberLock sync.Mutex
var sequenceNumber *uint64 = nil

func setSequenceNumber(s *uint64) {
	if s != nil {
		sequenceNumberLock.Lock()
		sequenceNumber = s
		sequenceNumberLock.Unlock()
	}
}
