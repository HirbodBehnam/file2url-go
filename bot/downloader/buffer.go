package downloader

import "sync"

// bufferPool is used to create a pool to decrease the allocations needed for creating
// a buffer in streaming
var bufferPool = sync.Pool{New: func() interface{} {
	return make([]byte, defaultPartSize)
}}

// toWriteBlock contains two byte slices:
type toWriteBlock struct {
	// toWrite contains the slice to be written in io.Writer
	toWrite []byte
	// original contains the slice which must be put in bufferPool after usage
	// toWrite is must be a slice of original
	original []byte
}
