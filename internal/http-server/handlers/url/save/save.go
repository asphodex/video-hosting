package save

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"video-hosting/internal/lib/random"
)

type VideoSaver interface {
	SaveVideo(url string, videoName string, author string) (int64, error)
}

func New(log *slog.Logger, videoSaver VideoSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		err := r.ParseMultipartForm(10 << 20) // max 10mb
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unable to parse multipart form", slog.String("error", err.Error()))
			return
		}
		author := r.Form.Get("author")
		if author == "" {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unable to parse author param", slog.String("error", err.Error()))
			return
		}
		f, handler, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unable to parse file", slog.String("error", err.Error()))
			return
		}
		defer func(f multipart.File) {
			err := f.Close()
			if err != nil {
				log.Error("Unable to close multipart file", slog.String("error", err.Error()))
				return
			}
		}(f)

		url := random.NewRandomString(10)

		name := r.Form.Get("name")
		if name == "" {
			name = "Unnamed video"
		}

		fileExtension := strings.ToLower(filepath.Ext(handler.Filename))
		if fileExtension != ".mp4" {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unsupported file extension", slog.String("error", err.Error()))
		}

		path := filepath.Join(".", "storage/data/videos")
		_ = os.MkdirAll(path, os.ModePerm)
		fullPath := path + "/" + url + fileExtension

		file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unable to open file for writing", slog.String("error", err.Error()))
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Error("Unable to close multipart file", slog.String("error", err.Error()))
				return
			}
		}(file)

		_, err = io.Copy(file, f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode("something went wrong")
			return
		}

		_, err = videoSaver.SaveVideo(url, name, author)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Unable to save video", slog.String("error", err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("File uploaded successfully")
		log.Info("File uploaded successfully", slog.String("name", name), slog.String("path", path))
		return
	}
}
