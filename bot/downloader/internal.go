package downloader

import (
	"context"
	"fmt"
	"io"

	"github.com/gotd/td/tg"
)

type chunkSource struct {
	loc tg.InputFileLocationClass
	api *tg.Client
	// What byte is the end
	end int64
}

// Chunk implements ChunkSource.
func (s chunkSource) Chunk(ctx context.Context, offset int64, b []byte) (int64, error) {
	req := &tg.UploadGetFileRequest{
		Offset:   offset,
		Limit:    len(b),
		Location: s.loc,
	}
	req.SetPrecise(true)
	req.SetCDNSupported(true)

	r, err := s.api.UploadGetFile(ctx, req)
	if err != nil {
		return 0, err
	}

	switch result := r.(type) {
	case *tg.UploadFile:
		n := int64(copy(b, result.Bytes))

		var err error
		if int64(req.Offset)+n >= s.end {
			// No more data.
			err = io.EOF
		}

		return n, err
	default:
		return 0, fmt.Errorf("unexpected type %T", r)
	}
}
