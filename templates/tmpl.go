package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
)

var success = map[string]bool{}

var funcs = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
}

//go:embed built-in/*.tmpl
var templates embed.FS

// Execute runs built-in and user-defined templates and links user defined files
func Execute(color models.Output) {
	defer link()

	if err := os.MkdirAll(config.CacheDir, 755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app cache directory: %v\n", err)
		return
	}

	defaultTmpl := template.New("").Funcs(funcs)

	tmpls, err := defaultTmpl.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse templates: %v\n", err)
		return
	}

	for _, tmpl := range tmpls.Templates() {
		execute(tmpl, color)
	}

	templateRoot := filepath.Join(config.ConfigDir, "templates")

	templatePath, err := filepath.Glob(templateRoot + "/*.tmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed find template: %v\n", err)
		return
	}

	tmpls, err = defaultTmpl.ParseFiles(templatePath...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse templates: %v\n", err)
		return
	}

	for _, tmpl := range tmpls.Templates() {
		execute(tmpl, color)
	}
}

func link() {
	if config.Global.Link == nil {
		return
	}

	for src, paths := range config.Global.Link {
		if _, ok := success[src]; !ok {
			fmt.Fprintf(os.Stderr, "Skiping %s. %s doesn't exists\n", src, src)
			continue
		}
		if !success[src] {
			fmt.Fprintf(os.Stderr, "Skiping %s. %s failed\n", src, src)
			continue
		}

		for path := range slices.Values(paths) {
			path, err := config.FindPath(config.ConfigDir, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to find path: %v\n", err)
				return
			}

			srcDir := filepath.Join(config.CacheDir, src)
			err = hardlinkOrCopy(srcDir, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to HardlinkOrCopy: %v", err)
				continue
			}

			fmt.Fprintf(os.Stderr, "Linked %s -> %s", srcDir, path)
		}
	}
}

// execute executes a template using color
func execute(tmpl *template.Template, color models.Output) {
	name := tmpl.Name()

	saveFile := strings.TrimSuffix(name, ".tmpl")
	outputPath := filepath.Join(config.CacheDir, saveFile)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", outputPath, err)
		success[saveFile] = false
		return
	}

	err = tmpl.Execute(file, color)
	file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template %s: %v\n", name, err)
		success[saveFile] = false
		return
	}

	success[saveFile] = true
	fmt.Printf("Template %s written to %s\n", name, outputPath)
}
