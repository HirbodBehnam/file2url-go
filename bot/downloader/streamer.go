package downloader

import (
	"context"
	"file2url/util"
	"fmt"
	"io"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
)

// ChunkSource downloads chunks.
type ChunkSource interface {
	Chunk(ctx context.Context, offset int64, b []byte) (int64, error)
}

// Streamer provides a pseudo-stream.
type Streamer struct {
	source ChunkSource // source of chunks
}

// nearestOffset returns the nearest offset that will conform to aligning
// requirements.
func nearestOffset(align, offset int64) int64 {
	if align == 0 {
		return offset
	}
	if offset == 0 {
		return 0
	}
	return offset - (offset % align)
}

func (s Streamer) safeRead(ctx context.Context, offset int64, data []byte) (int64, error) {
	for {
		n, err := s.source.Chunk(ctx, offset, data)
		if flood, err := tgerr.FloodWait(ctx, err); err != nil {
			if flood || tgerr.Is(err, tg.ErrTimeout) {
				continue
			}
			return n, err
		}
		if n < 0 || n > int64(len(data)) {
			return n, fmt.Errorf("invalid chunk: %d", n)
		}

		return n, nil
	}
}

func checkDone(ctx context.Context, errorChan <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}

// StreamAt streams from reader to "w" with "skip" offset.
func (s Streamer) StreamAt(ctx context.Context, skip, toRead int64, w io.Writer) error {
	var (
		offset         = nearestOffset(defaultPartSize, skip)
		bufSkip        = skip - offset
		toWrite        = make(chan toWriteBlock, 1)
		writeErrorChan = make(chan error, 1)
		closer         = util.SingleThreadOnce{}
	)
	defer closer.Do(func() {
		close(toWrite)
	})
	// Create a goroutine for writing
	go writerGoroutine(toWrite, writeErrorChan, w)
	// Read in loop
	for {
		if err := checkDone(ctx, writeErrorChan); err != nil {
			return err
		}
		buf := bufferPool.Get().([]byte)
		nr, er := s.safeRead(ctx, offset, buf)
		if er != nil && er != io.EOF {
			// Reading side done with error
			bufferPool.Put(buf)
			return er
		}

		if nr > 0 {
			bufferSlice := buf[bufSkip:nr]
			if toRead < int64(len(bufferSlice)) {
				bufferSlice = bufferSlice[:toRead]
			}
			toWrite <- toWriteBlock{
				toWrite:  bufferSlice,
				original: buf,
			}
			toRead -= int64(len(bufferSlice))
		}
		if er == io.EOF || toRead <= 0 {
			// Reading side exhausted.
			closer.Do(func() {
				close(toWrite)
			})
			return <-writeErrorChan // wait until everything is written
		}

		// Continue.
		offset += defaultPartSize // next chunk
		bufSkip = 0               // only skip at first chunk
	}
}

// writerGoroutine writes whatever comes in toWrite into w
// Writes until toWrite is closed
// At the end, send either nil (on toWrite closed) or the error of w.Write() in errorChan
func writerGoroutine(toWrite <-chan toWriteBlock, errorChan chan<- error, w io.Writer) {
	for {
		data, ok := <-toWrite
		if !ok {
			errorChan <- nil
			return
		}
		_, err := w.Write(data.toWrite)
		bufferPool.Put(data.original) // put back the buffer in original pool
		if err != nil {
			errorChan <- err
			return
		}
	}
}

// NewStreamer initializes and returns new *Streamer using provided chunk
// source and chunk size.
func NewStreamer(r ChunkSource) *Streamer {
	return &Streamer{
		source: r,
	}
}
