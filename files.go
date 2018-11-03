package got

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// locate templates in possibly nested subfolders
func findTemplatesRecursively(path string, extension string) (paths []string, err error) {
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

// Handles reading templates files in the given directory + ending path
func loadTemplateFiles(dir, path, extension string) (templates map[string][]byte, err error) {
	var files []string
	files, err = findTemplatesRecursively(filepath.Join(dir, path), extension)
	if err != nil {
		return
	}

	templates = make(map[string][]byte)

	for _, path = range files {
		var b []byte
		b, err = ioutil.ReadFile(path)
		if err != nil {
			return
		}

		// Convert "templates/layouts/base.html" to "layouts/base"
		// For subfolders the extra folder name is included:
		// "templates/includes/sidebar/ad1.html" to "includes/sidebar/ad1"
		name := strings.TrimPrefix(filepath.Clean(path), filepath.Clean(dir)+"/")
		name = strings.TrimSuffix(name, filepath.Ext(name))

		templates[name] = b
	}

	return
}
