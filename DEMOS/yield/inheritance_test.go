package main

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

// Was trying to get yield to work, ultimatily gave up

var layout = `layout: {{.Name}}: {{template "content" .}} {{block "footer" .}}no{{end}}`
var sidebar = `sidebar: {{.Name}}`
var home = `{{define "content"}}home: {{template "sidebar/one" .}}{{end}}`
var about = `{{define "content"}}about: {{template "sidebar/one" .}}{{end}}`
var footer = `Footer`

var CorrectOutput = "layout: John: home: sidebar: John Footer"

func TestOrder(t *testing.T) {

	// files := map[string]string{
	// 	"home":        home,
	// 	"about":       about,
	// 	"sidebar/one": sidebar,
	// 	"layout/base": layout,
	// }

	templates := template.New("templates").Funcs(template.FuncMap{
		"yield": func() (string, error) {
			return "", fmt.Errorf("yield called unexpectedly.")
		},
	})

	template.Must(templates.New("sidebar/one").Parse(sidebar))
	template.Must(templates.New("layout/base").Parse(layout))
	template.Must(templates.New("footer").Parse(footer))
	template.Must(templates.New("home").Parse(home))

	// for name, templateString := range files {
	// 	template.Must(templates.New(name).Parse(templateString))
	// }

	fmt.Println(templates.DefinedTemplates())

	data := struct{ Name string }{Name: "John"}

	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, "layout/base", data)
			return template.HTML(buf.String()), err
		},
	}
	templates.Funcs(funcs)

	var b []byte
	buf := bytes.NewBuffer(b)

	// err := templates.Execute(buf, data)
	err := templates.ExecuteTemplate(buf, "layout/base", struct{ Name string }{Name: "John"})
	if err != nil {
		fmt.Println("buf", buf.Bytes())
		t.Error(err)
	}

	have := string(buf.Bytes())
	want := CorrectOutput

	if have != want {
		t.Errorf("Invalid Response:\n\tgot:  %q\n\twant: %q", have, want)
	}

}

//
// func TestDefault(t *testing.T) {
// 	tmpl := template.New("foo")
// 	tmpl.Parse(home)
// 	tmpl.Parse(layout)
// 	tmpl.Parse(sidebar)
//
// 	fmt.Println(tmpl.DefinedTemplates())
//
// 	var b []byte
// 	buf := bytes.NewBuffer(b)
//
// 	err := tmpl.Execute(buf, nil)
// 	// err := tmpl.ExecuteTemplate(os.Stdout, "layout", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	have := string(buf.Bytes())
// 	want := "layout: home: sidebar"
//
// 	if have != want {
// 		t.Errorf("Invalid Response:\n\tgot:  %q\n\twant: %q", have, want)
// 	}
//
// }

//
// func TestNested(t *testing.T) {
// 	tmpl := template.New("foo")
// 	tmpl.Parse(layout)
// 	tmpl.Parse(sidebar)
// 	tmpl.Parse(home)
//
// 	fmt.Println(tmpl.DefinedTemplates())
//
// 	var b []byte
// 	buf := bytes.NewBuffer(b)
//
// 	err := tmpl.Execute(buf, nil)
// 	// err := tmpl.ExecuteTemplate(os.Stdout, "layout", nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	have := string(buf.Bytes())
// 	want := "Layout: Home: Sidebar"
//
// 	if have != want {
// 		t.Errorf("Invalid Response:\n\tgot:  %q\n\twant: %q", have, want)
// 	}
//
// }
