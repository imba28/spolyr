package api

import (
	"net/http"
	"os"
	"path/filepath"
)

func spaFileHandler(publicFolder string) http.HandlerFunc {
	const indexPath = "index.html"
	staticHandler := http.FileServer(http.Dir(publicFolder))

	return func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		path = filepath.Join(publicFolder, path)

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(publicFolder, indexPath))
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		staticHandler.ServeHTTP(w, r)
	}
}
