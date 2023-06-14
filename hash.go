package routines

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func hashCaller(callerStack []uintptr) string {
	pcHash := make([]byte, 0, len(callerStack)*4)
	for _, caller := range callerStack {
		pcHash = append(pcHash, uIntPtrToBytes(caller)...)
	}
	return string(pcHash)
}

func uIntPtrToBytes(u uintptr) []byte {
	size := unsafe.Sizeof(u)

	b := make([]byte, size)
	switch size {
	case 4:
		binary.LittleEndian.PutUint32(b, uint32(u))
	case 8:
		binary.LittleEndian.PutUint64(b, uint64(u))
	default:
		panic(fmt.Errorf("unknown uintptr size: %d", size))
	}

	return b
}
