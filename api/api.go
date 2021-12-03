package api

import (
	"file2url/bot/downloader"
	"file2url/shared"
	"github.com/gorilla/mux"
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
	w.Header().Set("Content-Length", "10")
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(file.Name))
	// w.Header().Set("Accept-Ranges", "bytes")
	// Run another client
	err := downloader.Download(r.Context(), shared.API, &tg.InputDocumentFileLocation{
		ID:            file.ID,
		AccessHash:    file.AccessHash,
		FileReference: file.FileReference,
	}, w, 0, int64(file.Size))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("cannot start client for download: ", err)
		return
	}
}
