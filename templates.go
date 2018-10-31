package got

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

// Children define the base template using comments: { /* use basetemplate */ }
var parentRegex = regexp.MustCompile(`\{\s*\/\*\s*use\s(\w+)\s*\*\/\s*\}`)

// Error for loading missing templates
// var ErrNotFound = errors.New("Template not found")

// NotFoundError for type assertions in the caller while still providing context
type NotFoundError struct {
	Name string
}

func (t *NotFoundError) Error() string {
	return fmt.Sprintf("Template %q not found", t.Name)
}

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
// func FindTemplates(path string, extension string) (paths []string, err error) {
// 	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
// 		if err == nil {
// 			if strings.Contains(path, extension) {
// 				paths = append(paths, path)
// 			}
// 		}
// 		return err
// 	})
// 	return
// }

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

func LoadTemplateFiles(dir, path string) (templates map[string][]byte, err error) {
	var files []string
	files, err = filepath.Glob(filepath.Join(dir, path))
	if err != nil {
		return
	}

	templates = make(map[string][]byte)

	for _, path = range files {
		// fmt.Printf("Loading: %s\n", path)
		var b []byte
		b, err = ioutil.ReadFile(path)
		if err != nil {
			return
		}

		// Convert "templates/layouts/base.html" to "layouts/base"
		name := strings.TrimPrefix(filepath.Clean(path), filepath.Clean(dir)+"/")
		name = strings.TrimSuffix(name, filepath.Ext(name))

		// fmt.Printf("%q = %q\n", name, b)
		templates[name] = b
	}

	return
}

func (t *Templates) Load(templatesDir string) (err error) {

	// Child pages to render
	var pages map[string][]byte
	pages, err = LoadTemplateFiles(templatesDir, "pages/*"+t.Extension)
	if err != nil {
		return
	}

	// Shared templates across multiple pages (sidebars, scripts, footers, etc...)
	var includes map[string][]byte
	includes, err = LoadTemplateFiles(templatesDir, "includes/*"+t.Extension)
	if err != nil {
		return
	}

	// Layouts used by pages
	var layouts map[string][]byte
	layouts, err = LoadTemplateFiles(templatesDir, "layouts/*"+t.Extension)
	if err != nil {
		return
	}

	var tmpl *template.Template
	for name, b := range pages {

		matches := parentRegex.FindSubmatch(b)
		basename := filepath.Base(name)

		tmpl, err = template.New(basename).Parse(string(b))

		// Uses a layout
		if len(matches) == 2 {

			l, ok := layouts[filepath.Join("layouts", string(matches[1]))]
			if !ok {
				err = fmt.Errorf("Unknown layout %s%s\n", matches[1], t.Extension)
				return
			}

			tmpl.New("layout").Parse(string(l))
		}

		if len(includes) > 0 {
			for name, src := range includes {
				// fmt.Printf("\tAdding:%s = %s\n", name, string(src))
				_, err = tmpl.New(name).Parse(string(src))
				if err != nil {
					return
				}
			}
		}

		t.Templates[basename] = tmpl
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
// func (t *Templates) CompilePool(template string, data interface{}, buf *bytes.Buffer) error {
//
// 	// fmt.Println("Compile:", template)
//
// 	// Look for the template
// 	tmpl, ok := t.Templates[template]
//
// 	if !ok {
// 		err := &NotFoundError{template}
// 		return err
// 	}
//
// 	// fmt.Printf("\t%s\n", tmpl.DefinedTemplates())
//
// 	// Create a buffer so syntax errors don't return a half-rendered response body
// 	// buf := bufpool.Get()
// 	// defer bufpool.Put(buf) // TODO fix this as it removes the content!
//
// 	// if err := tmpl.Execute(buf, data); err != nil {
// 	if err := tmpl.ExecuteTemplate(buf, "layout", data); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// // Compile the template and return the buffer containing the rendered bytes
func (t *Templates) Compile(template string, data interface{}) (*bytes.Buffer, error) {

	// Look for the template
	tmpl, ok := t.Templates[template]

	if !ok {
		err := &NotFoundError{template}
		return nil, err
	}

	// fmt.Printf("\t%s\n", tmpl.DefinedTemplates())

	// Create a buffer so syntax errors don't return a half-rendered response body
	var b []byte
	buf := bytes.NewBuffer(b)

	if err := tmpl.ExecuteTemplate(buf, "layout", data); err != nil {
		return buf, err
	}

	return buf, nil
}

// Make sure any template errors are caught before sending content to client
// A BufferPool will reduce allocs
// var bufpool *bpool.BufferPool
// var bufpool *bpool.SizedBufferPool
//
// func init() {
// 	// bufpool = bpool.NewBufferPool(32)
// 	bufpool = bpool.NewSizedBufferPool(1024, 64)
// }
