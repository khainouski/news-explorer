// Package web parses and renders the html/template pages under web/ (layouts, components,
// pages), embedded into the binary so the running container needs no files on disk beyond the
// static assets it serves directly (web/static).
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

// components is a directory (not a *.html glob) so it embeds its whole subtree - components live
// in semantic subfolders (article/, source/, navigation/, shared/), not flat.
//
//go:embed layouts/*.html components pages/*.html
var templatesFS embed.FS

//go:embed static
var StaticFS embed.FS

// pages lists every file under web/pages/ by name (without extension) - add new ones here so
// they get parsed at startup.
var pages = []string{"home", "coming_soon", "sources", "source_form", "not_found", "login", "change_password"}

var parsed = map[string]*template.Template{}

func init() {
	componentFiles, err := componentTemplateFiles()
	if err != nil {
		panic(fmt.Sprintf("web: list component templates: %v", err))
	}

	for _, p := range pages {
		files := append([]string{"layouts/base.html"}, componentFiles...)
		files = append(files, "pages/"+p+".html")

		t, err := template.ParseFS(templatesFS, files...)
		if err != nil {
			panic(fmt.Sprintf("web: parse template %q: %v", p, err))
		}

		parsed[p] = t
	}
}

// componentTemplateFiles walks components/ for every .html file, however deeply nested under its
// subfolders - template.ParseFS's glob has no recursive wildcard, so this stands in for one.
func componentTemplateFiles() ([]string, error) {
	var files []string

	err := fs.WalkDir(templatesFS, "components", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// Render executes the named page (see pages above) against data, writing HTML to w.
func Render(w http.ResponseWriter, page string, data any) error {
	return RenderBlock(w, page, "base", data)
}

// RenderBlock executes one named {{define}} block of page instead of the full "base" layout -
// for HTMX requests that swap a fragment of the page rather than navigating to it.
func RenderBlock(w http.ResponseWriter, page, block string, data any) error {
	t, ok := parsed[page]
	if !ok {
		return fmt.Errorf("web: unknown page %q", page)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	return t.ExecuteTemplate(w, block, data)
}
