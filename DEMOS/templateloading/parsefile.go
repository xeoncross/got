package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	//"io/ioutil"
	"strings"
)

func main() {

	t := template.New("foo")
	t.Parse(`{{define "content"}}foo contents{{end}}`)

	var err error
	_, err = ParseFile(t, "sub")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(os.Stdout, nil)

	fmt.Println(t)
}

// ParseFile with custom name instead of using .ParseFiles()
func ParseFile(t *template.Template, path string) (tmpl *template.Template, err error) {
	basename := filepath.Base(path)
	basename = strings.TrimSuffix(basename, filepath.Ext(basename))

	fmt.Println(path, basename)

	var b []byte
	/*
		b, err = ioutil.ReadFile(path)
		if err != nil {
			return
		}
	*/

	b = []byte(`Sub Template: {{template "content"}}`)

	// tmpl = t.New(basename)
	_, err = t.Parse(string(b))

	if err != nil {
		return
	}

	// fmt.Println(tmpl.DefinedTemplates())
	fmt.Println(t.DefinedTemplates())
	return
}
