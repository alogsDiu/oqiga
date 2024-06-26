package renderer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
)

type RendererMachine struct {
	Templates *template.Template
}

var Renderer RendererMachine

func CreateRenderer(pattern string) {
	Renderer = RendererMachine{Templates: loadTemplates(pattern)}
}

func loadTemplates(pattern string) *template.Template {
	templates := template.New("")
	layouts, err := filepath.Glob(pattern + "/layouts/*.html")

	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(pattern + "/block/*.html")

	if err != nil {
		log.Fatal(err)
	}

	pages, err := filepath.Glob(pattern + "/pages/*.html")

	if err != nil {
		log.Fatal(err)
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		files = append(files, pages...)
		templates = template.Must(templates.ParseFiles(files...))
	}

	return templates
}

func (re *RendererMachine) Render(w http.ResponseWriter, name string, r *http.Request, data interface{}) {
	err := re.Templates.ExecuteTemplate(w, name, data)

	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}

	res, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(res))
}
