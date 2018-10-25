package main

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

var layout = `layout: {{.Name}}: {{block "content" .}}default{{end}} {{block "footer" .}}no{{end}}`
var sidebar = `{{define "sidebar"}}sidebar: {{.Name}}{{end}}`
var home = `{{define "content"}}home: {{template "sidebar" .}}{{end}}`
var footer = `{{block "content"}}home_footer{{end}}{{define "footer"}}Footer{{end}}`

var CorrectOutput = "layout: John: home: sidebar: John Footer"

func TestOrder(t *testing.T) {

	sets := map[string][]string{
		// Bad
		// "home,layout,sidebar": []string{home, layout, sidebar, footer},
		// "home,sidebar,layout": []string{home, sidebar, layout, footer},
		// "sidebar,home,layout": []string{sidebar, home, layout, footer},
		// Good
		"layout,home,sidebar": []string{layout, footer, home, sidebar},
		"layout,sidebar,home": []string{layout, sidebar, home, footer},
		"sidebar,layout,home": []string{sidebar, layout, home, footer},
	}

	for name, templates := range sets {

		t.Run(name, func(t *testing.T) {
			tmpl := template.New("foo")
			for _, str := range templates {
				tmpl.Parse(str)
			}
			fmt.Println(tmpl.DefinedTemplates())

			var b []byte
			buf := bytes.NewBuffer(b)

			err := tmpl.Execute(buf, struct{ Name string }{Name: "John"})
			// err := tmpl.ExecuteTemplate(buf, "layout", nil)
			if err != nil {
				t.Error(err)
			}

			have := string(buf.Bytes())
			want := CorrectOutput

			if have != want {
				t.Errorf("Invalid Response:\n\tgot:  %q\n\twant: %q", have, want)
			}
		})
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
