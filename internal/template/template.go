package template

import (
	"github.com/gin-gonic/gin"
	"github.com/imba28/spolyr/internal/template/files"
	"net/http"
	"strings"
)
import "html/template"

var (
	homepageTemplate        = parse("pages/index.html")
	trackTemplate           = parse("pages/track.html")
	trackEditTemplate       = parse("pages/track-edit.html")
	trackLyricsSyncTemplate = parse("pages/track-lyrics-sync.html")
	tracksTemplate          = parse("pages/tracks.html")
	errorTemplate           = parse("pages/error.html")
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
	return executeWithStatus(homepageTemplate, w, p, status)
}
func TrackPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(trackTemplate, w, p, status)
}
func TrackEditPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(trackEditTemplate, w, p, status)
}
func TrackLyricsSyncPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(trackLyricsSyncTemplate, w, p, status)
}
func TracksPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(tracksTemplate, w, p, status)
}
func ErrorPage(w http.ResponseWriter, p gin.H, status int) error {
	return executeWithStatus(errorTemplate, w, p, status)
}
