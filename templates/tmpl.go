package templates

import (
	"embed"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Nadim147c/rong/internal/config"
	"github.com/Nadim147c/rong/internal/models"
)

var success = map[string]bool{}

//go:embed built-in/*.tmpl
var templates embed.FS

// Execute runs built-in and user-defined templates and links user defined files
func Execute(color models.Output) {
	defer link()

	if err := os.MkdirAll(config.CacheDir, 755); err != nil {
		slog.Error("Failed to create app cache directory", "error", err)
		return
	}

	defaultTmpl := template.New("").Funcs(funcs)

	tmpls, err := defaultTmpl.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return
	}

	for _, tmpl := range tmpls.Templates() {
		execute(tmpl, color)
	}

	templateRoot := filepath.Join(config.ConfigDir, "templates")

	templatePath, err := filepath.Glob(templateRoot + "/*.tmpl")
	if err != nil {
		slog.Error("Failed to find template", "error", err)
		return
	}

	if len(templatePath) == 0 {
		return
	}

	tmpls, err = defaultTmpl.ParseFiles(templatePath...)
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
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
			slog.Warn("Skipping source, it doesn't exist", "src", src)
			continue
		}
		if !success[src] {
			slog.Warn("Skipping source, it previously failed", "src", src)
			continue
		}

		for path := range slices.Values(paths) {
			path, err := config.FindPath(config.ConfigDir, path)
			if err != nil {
				slog.Error("Failed to find path", "error", err)
				return
			}

			srcDir := filepath.Join(config.CacheDir, src)
			hardlinkOrCopy(srcDir, path)
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
		slog.Error("Error creating file", "file", outputPath, "error", err)
		success[saveFile] = false
		return
	}

	err = tmpl.Execute(file, color)
	file.Close()
	if err != nil {
		slog.Error("Error executing template", "template", name, "error", err)
		success[saveFile] = false
		return
	}

	success[saveFile] = true
	slog.Info("Template written", "template", name, "path", outputPath)
}
