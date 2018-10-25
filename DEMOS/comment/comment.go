package comment

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// Children define the base template using comments: { /* use basetemplate */ }
var parentRegex = regexp.MustCompile(`\{\s*\/\*\s*use\s(\w+)\s*\*\/\s*\}`)

func LoadTemplateFile(path string) (t *template.Template, err error) {
	var b []byte
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	t, err = template.New("").Parse(string(b))
	return
}

func LoadTemplates(dir, path string) (templates map[string]*template.Template, err error) {
	var files []string
	files, err = filepath.Glob(path)
	if err != nil {
		return
	}

	templates = make(map[string]*template.Template)

	var t *template.Template
	for _, path = range files {
		t, err = LoadTemplateFile(path)
		if err != nil {
			return
		}

		// Convert "templates/layouts/base.html" to "layouts/base"
		name := strings.TrimPrefix(path, dir)
		name = strings.TrimSuffix(name, filepath.Ext(name))

		templates[name] = t
	}

	return
}

func Load(templatesDir, extension string) (templates *template.Template, err error) {

	// Child pages to render
	var pages map[string]*template.Template
	pages, err = LoadTemplates(templatesDir, "pages/*"+extension)
	if err != nil {
		return
	}

	// Shared templates across multiple pages (sidebars, scripts, footers, etc...)
	var includes map[string]*template.Template
	includes, err = LoadTemplates(templatesDir, "includes/*"+extension)
	if err != nil {
		return
	}

	// // Layouts used by pages
	var layouts map[string]*template.Template
	layouts, err = LoadTemplates(templatesDir, "layouts/*"+extension)
	if err != nil {
		return
	}

	_ = includes
	_ = layouts

	// var b []byte
	for name, tmpl := range pages {
		// 	b, err = ioutil.ReadFile(page)
		// 	if err != nil {
		// 		return
		// 	}
		//
		// 	matches := parentRegex.FindSubmatch(b)
		//
		// 	// Does not use a layout
		// 	if len(matches) == 0 {
		// 		fmt.Printf("No layout for %s\n", page)
		// 	} else {
		// 		fmt.Printf("%q\n", matches)
		// 	}
		//
		// 	// name := strings.TrimPrefix(tm, trimPrefix)
		// 	// template.Must(templates.New(name).Parse(string(b)))
	}

	return
}
