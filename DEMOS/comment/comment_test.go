package comment

import (
	"bytes"
	"fmt"
	"testing"
)

/* Golang requires loading child templates from the parent. This is backwards
when trying to write a page and specify (from that page) which layout to use.
Using a template comment, we can have a page define the layout it inherits
from. */

var CorrectOutput = "layout: John: home: sidebar: John Footer"

func TestOrder(t *testing.T) {

	templates, err := Load("templates", ".html")

	if err != nil {
		t.Error(err)
	}

	for name, tmpl := range templates {
		fmt.Printf("\t%s = %s\n", name, tmpl.DefinedTemplates())

		var b []byte
		buf := bytes.NewBuffer(b)

		// err := templates.Execute(buf, data)
		err := tmpl.ExecuteTemplate(buf, "home", struct{ Name string }{Name: "John"})
		if err != nil {
			// fmt.Println("buf", buf.Bytes())
			t.Error(err)
		}

		have := string(buf.Bytes())
		want := CorrectOutput

		if have != want {
			t.Errorf("Invalid Response:\n\tgot:  %q\n\twant: %q", have, want)
		}
	}

}
