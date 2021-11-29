package api

import (
	"context"
	"file2url/shared"
	"github.com/gorilla/mux"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"log"
	"net/http"
	"strconv"
)

// downloadEndpoint downloads a file requested
func downloadEndpoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	file, exists := shared.Database.Load(id)
	if !exists {
		http.NotFound(w, r)
		return
	}
	// Set the headers
	w.Header().Set("Content-Length", strconv.Itoa(file.Size))
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(file.Name))
	// w.Header().Set("Accept-Ranges", "bytes")
	// Run another client
	err := telegram.BotFromEnvironment(r.Context(), telegram.Options{}, noOpClient, func(ctx context.Context, client *telegram.Client) error {
		_, err := downloader.NewDownloader().Download(client.API(), &tg.InputDocumentFileLocation{
			ID:            file.ID,
			AccessHash:    file.AccessHash,
			FileReference: file.FileReference,
		}).Stream(ctx, w)
		return err
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("cannot start client for download: ", err)
		return
	}
}

func noOpClient(context.Context, *telegram.Client) error {
	return nil
}
