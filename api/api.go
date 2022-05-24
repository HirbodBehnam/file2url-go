package api

import (
	"file2url/bot/clients"
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
	// Filter the methods
	if r.Method != "GET" && r.Method != "HEAD" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get the file id
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
	// Check head method
	if r.Method == "HEAD" {
		return
	}
	// Run another client
	clientAPI, err := clients.ClientPools.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		log.Println("cannot get client from pool:", err)
		return
	}
	defer clients.ClientPools.Put(clientAPI)
	err = downloader.Download(r.Context(), clientAPI.API(), &tg.InputDocumentFileLocation{
		ID:            file.ID,
		AccessHash:    file.AccessHash,
		FileReference: file.FileReference,
	}, w, from, to+1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("cannot start client for download: ", err, r.Header.Get("Range"))
		return
	}
}
