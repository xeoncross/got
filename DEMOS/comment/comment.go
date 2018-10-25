package comment

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// Children define the base template using comments: { /* use basetemplate */ }
var parentRegex = regexp.MustCompile(`\{\s*\/\*\s*use\s(\w+)\s*\*\/\s*\}`)

// func LoadTemplateFile(path string) (t *template.Template, err error) {
// 	var b []byte
// 	b, err = ioutil.ReadFile(path)
// 	if err != nil {
// 		return
// 	}
//
// 	t, err = template.New("").Parse(string(b))
// 	return
// }

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

		// fmt.Printf("%q\n", b)

		// Convert "templates/layouts/base.html" to "layouts/base"
		name := strings.TrimPrefix(filepath.Clean(path), filepath.Clean(dir)+"/")
		name = strings.TrimSuffix(name, filepath.Ext(name))

		templates[name] = b
	}

	return
}

func Load(templatesDir, extension string) (templates map[string]*template.Template, err error) {

	// Child pages to render
	var pages map[string][]byte
	pages, err = LoadTemplateFiles(templatesDir, "pages/*"+extension)
	if err != nil {
		return
	}

	// Shared templates across multiple pages (sidebars, scripts, footers, etc...)
	var includes map[string][]byte
	includes, err = LoadTemplateFiles(templatesDir, "includes/*"+extension)
	if err != nil {
		return
	}

	// Layouts used by pages
	var layouts map[string][]byte
	layouts, err = LoadTemplateFiles(templatesDir, "layouts/*"+extension)
	if err != nil {
		return
	}

	// Get ready to populate
	templates = make(map[string]*template.Template)

	// _ = includes
	_ = layouts

	var t *template.Template
	for name, b := range pages {

		fmt.Println(name)

		matches := parentRegex.FindSubmatch(b)

		t, err = template.New(name).Parse(string(b))

		// Uses a layout
		if len(matches) == 2 {

			l, ok := layouts[filepath.Join("layouts", string(matches[1]))]
			if !ok {
				err = fmt.Errorf("Unknown layout %s%s\n", matches[1], extension)
				return
			}

			t.New("layout").Parse(string(l))
		}

		if len(includes) > 0 {
			for name, src := range includes {
				// fmt.Printf("\tAdding:%s\n", name)
				_, err = t.New(name).Parse(string(src))
				if err != nil {
					return
				}
			}
		}

		templates[name] = t
	}

	return
}
