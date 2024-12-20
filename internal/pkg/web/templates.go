package web

import (
	"html/template"
	"io/fs"
	"os"
	"strings"
)

func TemplateParseFSRecursive(
	templates fs.FS,
	templatesDir string,
	ext string,
	funcMap template.FuncMap) (*template.Template, error) {

	templatesDirParts := strings.Split(templatesDir, string(os.PathSeparator))
	templatesDirPartsNum := len(templatesDirParts)

	root := template.New("")
	err := fs.WalkDir(templates, templatesDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}
			b, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}
			parts := strings.Split(path, string(os.PathSeparator))
			name := strings.Join(parts[templatesDirPartsNum:], string(os.PathSeparator))
			t := root.New(name).Funcs(funcMap)
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return root, err
}
