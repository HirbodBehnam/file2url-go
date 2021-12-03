package downloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"io"
)

// ChunkSource downloads chunks.
type ChunkSource interface {
	Chunk(ctx context.Context, offset int64, b []byte) (int64, error)
}

// Streamer provides a pseudo-stream.
type Streamer struct {
	align  int64       // required chunk size
	source ChunkSource // source of chunks
}

// nearestOffset returns nearest offset that will conform to aligning
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

// errInvalidWrite means that a write returned an impossible count.
var errInvalidWrite = errors.New("invalid write result")

func checkDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s Streamer) writeFull(ctx context.Context, buf []byte, dst io.Writer) (written int64, err error) {
	nr := len(buf)

	for {
		if err = checkDone(ctx); err != nil {
			break
		}
		if written == int64(nr) {
			break
		}
		// Same logic as in io.Copy.
		nw, ew := dst.Write(buf[written:nr])
		if nw < 0 || nr < nw {
			nw = 0
			if ew == nil {
				ew = errInvalidWrite
			}
		}
		written += int64(nw)
		if ew != nil {
			err = ew
			break
		}
		if nr != nw {
			err = io.ErrShortWrite
			break
		}
	}

	return written, err
}

// StreamAt streams from reader to "w" with "skip" offset.
func (s Streamer) StreamAt(ctx context.Context, skip, toRead int64, w io.Writer) error {
	var (
		buf     = make([]byte, s.align)
		offset  = nearestOffset(s.align, skip)
		bufSkip = skip - offset
	)
	for {
		if err := checkDone(ctx); err != nil {
			return err
		}
		nr, er := s.safeRead(ctx, offset, buf)
		if er != nil && er != io.EOF {
			// Reading side done.
			return er
		}

		nr = min(nr, toRead+bufSkip)
		if nr > 0 {
			if _, err := s.writeFull(ctx, buf[bufSkip:nr], w); err != nil {
				// Writing side done.
				return err
			}
		}
		toRead -= nr - bufSkip
		if er == io.EOF || toRead <= 0 {
			// Reading side exhausted.
			return nil
		}

		// Continue.
		offset += s.align // next chunk
		bufSkip = 0       // only skip at first chunk
	}
}

// NewStreamer initializes and returns new *Streamer using provided chunk
// source and chunk size.
func NewStreamer(r ChunkSource, chunkSize int64) *Streamer {
	if chunkSize <= 0 {
		panic("invalid chunk size")
	}
	return &Streamer{
		align:  chunkSize,
		source: r,
	}
}

func min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
