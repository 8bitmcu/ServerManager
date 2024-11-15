package main

import (
	"html/template"
	"io"
	"strings"

	"github.com/jessevdk/go-assets"
)

// Load templates from the embedded assets file
func LoadTemplate(t *template.Template, ext string) error {
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ext) {
			continue
		}
		h, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return err
		}
	}
	return nil
}

// Find a file in the embedded assets file
func FindFile(filePath string) *assets.File {
	for _, file := range Assets.Files {
		if file.Path == filePath {
			return file
		}
	}
	return nil
}
