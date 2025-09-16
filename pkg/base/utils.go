package base

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// ConvertToString converts and connects all objects to a single string
func ConvertToString(data ...interface{}) string {
	converted := make([]string, len(data))
	for i, str := range data {
		converted[i] = fmt.Sprintf("%v", str)
	}
	return strings.Join(converted, "")
}

var entityIdCounter atomic.Int32

// NextEntityId provides next empty entity ID
func NextEntityId() int32 {
	// if the value is max signed, go back to zero
	entityIdCounter.CompareAndSwap(int32((^uint32(0))>>1), 0)
	return entityIdCounter.Add(1)
}
