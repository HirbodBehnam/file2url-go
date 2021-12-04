package api

import (
	"file2url/bot/downloader"
	"file2url/shared"
	"fmt"
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
	// Get the range if needed
	from, to, err := getRange(r.Header.Get("Range"), file.Size)
	if err != nil {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// Set the headers
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(file.Name))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(to-from+1, 10))
	if r.Header.Get("Range") != "" {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", from, to, file.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	// Run another client
	err = downloader.Download(r.Context(), shared.API, &tg.InputDocumentFileLocation{
		ID:            file.ID,
		AccessHash:    file.AccessHash,
		FileReference: file.FileReference,
	}, w, from, to+1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("cannot start client for download: ", err)
		return
	}
}
