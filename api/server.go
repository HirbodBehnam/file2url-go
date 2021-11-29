package api

import (
	"context"
	"file2url/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// StartWebserver starts the http webserver to serve the files
func StartWebserver(stopChannel <-chan struct{}) {
	r := mux.NewRouter()
	r.HandleFunc("/{id}/", downloadEndpoint)
	srv := &http.Server{
		Handler:     r,
		Addr:        config.Config.ListenAddress,
		ReadTimeout: 10 * time.Second,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	<-stopChannel
	_ = srv.Shutdown(context.Background())
}
