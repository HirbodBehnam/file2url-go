package downloader

import (
	"context"
	"io"

	"github.com/gotd/td/tg"
)

const defaultPartSize = 64 * 1024

// Download simply downloads a file from telegram
func Download(ctx context.Context, client *tg.Client, location tg.InputFileLocationClass, output io.Writer, from, to int64) error {
	streamer := NewStreamer(chunkSource{
		loc: location,
		api: client,
		end: to,
	})
	return streamer.StreamAt(ctx, from, to-from, output)
}
