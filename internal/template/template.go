package template

import (
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/template/files"
	"net/http"
	"strings"
)
import "html/template"

var (
	homepage        = parse("pages/index.html")
	search          = parse("pages/search.html")
	track           = parse("pages/track.html")
	trackEdit       = parse("pages/track-edit.html")
	trackLyricsSync = parse("pages/track-lyrics-sync.html")
	tracks          = parse("pages/tracks.html")
)

var templateFunctions = template.FuncMap{
	"formatHTML": func(raw string) template.HTML {
		return template.HTML(strings.ReplaceAll(raw, "\n", "<br/>"))
	},
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").
			Funcs(templateFunctions).
			ParseFS(files.TemplateFiles,
				"includes/navbar.html",
				"includes/footer.html",
				"includes/track-list.html",
				"layout.html",
				file))
}

func executeWithStatus(t *template.Template, w http.ResponseWriter, p gin.H, status int) error {
	w.WriteHeader(status)
	return t.Execute(w, p)
}

func HomePage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(homepage, w, p, status)
}
func SearchPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(search, w, p, status)
}
func TrackPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(track, w, p, status)
}
func TrackEditPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(trackEdit, w, p, status)
}
func TrackLyricsSyncPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(trackLyricsSync, w, p, status)
}
func TracksPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(tracks, w, p, status)
}
