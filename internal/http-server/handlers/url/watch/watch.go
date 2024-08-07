package watch

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type VideoWatcher interface {
	GetPath(url string) (string, error)
}

func New(log *slog.Logger, videoWatcher VideoWatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.watch.New"
		url := chi.URLParam(r, "url")
		if url == "" {
			log.Info("url is empty")
			http.Error(w, "url is empty", http.StatusBadRequest)
			return
		}
		path, err := videoWatcher.GetPath(url)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "video/mp4")
		fullFilePath := "storage/data/videos/" + path + ".mp4"
		http.ServeFile(w, r, fullFilePath)
	}
}