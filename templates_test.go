package got

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Borrowed from text/template
// https://golang.org/src/text/template/examplefiles_test.go

// templateFile defines the contents of a template to be stored in a file, for testing.
type templateFile struct {
	name     string
	contents string
}

var testingTemplateFiles = []templateFile{
	// We have two pages each using a different parent layout
	{"pages/home.html", `{{define "content"}}home {{.Name}}{{end}} {{/* use one */}}`},
	{"pages/about.html", `{{define "content"}}about {{.Name}}{{end}}{{/* use two */}}`},
	// We have two different layouts (using two different styles)
	{"layouts/one.html", `Layout 1: {{.Name}} {{block "content" .}}{{end}} {{block "includes/sidebar" .}}{{end}}`},
	{"layouts/two.html", `Layout 2: {{.Name}} {{template "content" .}} {{template "includes/sidebar" .}}`},
	// We have two includes shared among the pages
	{"includes/header.html", `header`},
	{"includes/sidebar.html", `sidebar {{.Name}}`},
}

func createTestDir(files []templateFile) (dir string, err error) {
	dir, err = ioutil.TempDir("", "template")
	if err != nil {
		return
	}
	for _, file := range files {

		// Create sub directory of file (if needed)
		fd := filepath.Dir(filepath.Join(dir, file.name))
		err = os.MkdirAll(fd, os.ModePerm)
		if err != nil {
			return
		}

		var f *os.File
		f, err = os.Create(filepath.Join(dir, file.name))
		if err != nil {
			return
		}
		defer f.Close()
		_, err = io.WriteString(f, file.contents)
		if err != nil {
			return
		}
	}
	return
}

//
// Tests
//

// Here we demonstrate loading a set of templates from a directory.
func TestTemplates(t *testing.T) {
	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		t.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer os.RemoveAll(dir)

	templates := New(".html")
	err = templates.Load(dir)

	if err != nil {
		t.Error(err)
	}

	//
	// 1: Test page + include + layout
	//
	data := struct{ Name string }{"John"}
	b, err := templates.Compile("home", data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	got := string(b.Bytes())
	want := "Layout 1: John home John sidebar John"

	if got != want {
		t.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
	}

	//
	// 2: Test layout isolation
	//
	data = struct{ Name string }{"Jane"}
	b, err = templates.Compile("about", data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	got = string(b.Bytes())
	want = "Layout 2: Jane about Jane sidebar Jane"

	if got != want {
		t.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
	}

	//
	// 3: Test HTTP handler
	//
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data = struct{ Name string }{"Bob"}
		err := templates.Render(w, "home", data, http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	router.ServeHTTP(rr, req)

	got = rr.Body.String()
	want = "Layout 1: Bob home Bob sidebar Bob"

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Error(rr.Body.String())
	}

	if got != want {
		t.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
	}
}

func BenchmarkCompile(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		b.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer os.RemoveAll(dir)

	templates := New(".html")
	err = templates.Load(dir)

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		body, err := templates.Compile("home", data)
		if err != nil {
			b.Error(err)
		}

		got := string(body.Bytes())
		want := "Layout 1: John home John sidebar John"

		if got != want {
			b.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
		}
	}
}

/*
func BenchmarkCompileBuffer(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		b.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer os.RemoveAll(dir)

	templates := New(".html")
	err = templates.Load(dir)

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		buf := bufpool.Get()

		err := templates.CompileWithBuffer("home", data, buf)
		if err != nil {
			b.Error(err)
		}

		got := string(buf.Bytes())
		want := "Layout 1: John home John sidebar John"

		bufpool.Put(buf)

		if got != want {
			b.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
		}
	}
}
*/

func BenchmarkNativeTemplates(b *testing.B) {

	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	dir, err := createTestDir(testingTemplateFiles)

	if err != nil {
		b.Error(err)
	}

	// Clean up after the test; another quirk of running as an example.
	defer os.RemoveAll(dir)

	t := template.New("")

	var by []byte
	for _, name := range []string{"pages/home", "layouts/one", "includes/sidebar"} {
		by, err = ioutil.ReadFile(filepath.Join(dir, name+".html"))
		if err != nil {
			b.Error(err)
		}
		_, err = t.New(name).Parse(string(by))
		if err != nil {
			b.Error(err)
		}
	}

	b.ResetTimer()

	data := struct{ Name string }{"John"}

	for i := 0; i < b.N; i++ {

		by = nil
		buf := bytes.NewBuffer(by)

		if err := t.ExecuteTemplate(buf, "layouts/one", data); err != nil {
			b.Error(err)
		}

		got := string(buf.Bytes())
		want := "Layout 1: John home John sidebar John"

		if got != want {
			b.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
		}
	}
}

/*
func TestTemplates(t *testing.T) {
	// Here we create a temporary directory and populate it with our sample
	// template definition files; usually the template files would already
	// exist in some location known to the program.
	// dir := createTestDir([]templateFile{
	// 	// T0.tmpl is a plain template file that just invokes T1.
	// 	{"T0.tmpl", `T0 invokes T1: ({{template "T1"}})`},
	// 	// T1.tmpl defines a template, T1 that invokes T2.
	// 	{"T1.tmpl", `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`},
	// 	// T2.tmpl defines a template T2.
	// 	{"T2.tmpl", `{{define "T2"}}This is T2{{end}}`},
	// })
	// // Clean up after the test; another quirk of running as an example.
	// defer os.RemoveAll(dir)

	// templates := New(".tmpl")
	templates := New(".html")
	err := templates.Load("samples/native/pages", "samples/native/layouts", "samples/native/includes")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates.Render(w, "home", nil, http.StatusOK)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
		}
	})
	router.ServeHTTP(rr, req)

	got := rr.Body.String()
	want := `T0 invokes T1: (T1 invokes T2: (This is T2))` + "\n"

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Error(rr.Body.String())
	}

	if got != want {
		t.Errorf("handler returned wrong body:\n\tgot:  %q\n\twant: %q", got, want)
	}

}
*/
