package got

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/oxtoacart/bpool"
)

// Might provide a default error template, probably not
// const errorTemplateHTML = `
// <!DOCTYPE html>
// <html>
// 	<head>
// 		<meta charset="UTF-8">
// 		<title>{{.Error}}</title>
// 		<style type="text/css">
// 			body { margin: 0 auto; max-width: 600px; }
// 			pre { background: #eee; padding: 1em; margin: 1em 0; }
// 		</pre>
// 	</head>
// 	<body>
// 		<h2>Error: {{.Error}}</h2>
// 		<pre>{{.Trace}}</pre>
// 	</body>
// </html>`
//
// // Template for displaying errors
// var ErrorTemplate = template.Must(template.New("error").Parse(errorTemplateHTML))

// FindTemplates in path recursively
func FindTemplates(path string, extension string) (paths []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if strings.Contains(path, extension) {
				paths = append(paths, path)
			}
		}
		return err
	})
	return
}

// ParseFile with custom name instead of using .ParseFiles()
func ParseFile(t *template.Template, path string) (err error) {
	// basename := filepath.Base(path)
	// basename = strings.TrimSuffix(basename, filepath.Ext(basename))

	// fmt.Println(path, basename)

	var b []byte
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	// tmpl := t.New(basename)
	// _, err = tmpl.Parse(string(b))

	_, err = t.Parse(string(b))

	if err != nil {
		return
	}

	// fmt.Println(path, t.DefinedTemplates())
	// t = tmpl
	return
}

// AddTemplates found in path recursively
func AddTemplates(Templates *template.Template, path string, extension string) (err error) {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if strings.Contains(path, extension) {

				fmt.Printf("\tAdding: %s\n", filepath.Base(path))
				// var Templates *template.Template

				// Templates named "filename.html"
				_, err = Templates.ParseFiles(path)

				// err = ParseFile(Templates, path)

				if err != nil {
					return err
				}
				// Templates = tmpl

			}
		}
		return err
	})
}

// Templates Collection
type Templates struct {
	Extension string
	Templates map[string]*template.Template
	Functions template.FuncMap
}

// DefaultFunctions for templates
var DefaultFunctions = template.FuncMap{
	// Allow unsafe injection into HTML
	"noescape": func(a ...interface{}) template.HTML {
		return template.HTML(fmt.Sprint(a...))
	},
	"title": strings.Title,
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"trim":  strings.TrimSpace,
}

// New Templates collection
func New(extension string) *Templates {
	return &Templates{
		Extension: extension,
		Templates: make(map[string]*template.Template),
		Functions: DefaultFunctions,
	}
}

// Load templates segregated by template name
// First item is the child pages, the next and following are the layouts/includes:
// Load("pages/", "layouts/", "partials/")
func (t *Templates) Load(paths ...string) (err error) {

	var pagesPath string
	pagesPath, paths = strings.TrimRight(paths[0], "/"), paths[1:]

	var pages []string
	pages, err = FindTemplates(pagesPath, t.Extension)
	if err != nil {
		return
	}

	for _, pagePath := range pages {
		basename := filepath.Base(pagePath)
		basename = strings.TrimSuffix(basename, filepath.Ext(basename))

		fmt.Println(pagePath, basename)

		// Load this template
		tmp := template.New("foobar").Funcs(t.Functions)
		// err = ParseFile(tmp, pagePath)
		// if err != nil {
		// 	return
		// }
		t.Templates[basename] = tmp
		// t.Templates[basename] = template.Must(template.New(basename).Funcs(t.Functions).ParseFiles(pagePath))

		// Each add all the includes, partials, and layouts
		if len(paths) > 0 {
			for _, templateDir := range paths {
				AddTemplates(t.Templates[basename], templateDir, t.Extension)
			}
		}

		// _, err = tmp.ParseFiles(pagePath)
		err = ParseFile(tmp, pagePath)
		if err != nil {
			return
		}

		fmt.Println(basename, t.Templates[basename].DefinedTemplates())

	}

	return
}

// Render the template & data to the ResponseWriter safely
func (t *Templates) Render(w http.ResponseWriter, template string, data interface{}, status int) error {

	buf, err := t.Compile(template, data)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)

	return nil
}

// Compile the template and return the buffer containing the rendered bytes
func (t *Templates) Compile(template string, data interface{}) (*bytes.Buffer, error) {

	fmt.Println("Complie:", template)

	// Look for the template
	tmpl, ok := t.Templates[template]

	if !ok {
		return nil, ErrNotFound
	}

	fmt.Printf("\t%s\n", tmpl.DefinedTemplates())

	// Create a buffer so syntax errors don't return a half-rendered response body
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		// if err := tmpl.ExecuteTemplate(buf, "content", data); err != nil {
		return buf, err
	}

	return buf, nil
}

// Error for loading missing templates
var ErrNotFound = errors.New("Template not found")

// Make sure any template errors are caught before sending content to client
// A BufferPool will reduce allocs
var bufpool *bpool.BufferPool

func init() {
	bufpool = bpool.NewBufferPool(64)
}
