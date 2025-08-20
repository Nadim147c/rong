package templates

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Nadim147c/go-config"
	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/pathutil"
)

var success = map[string]bool{}

//go:embed built-in/*.tmpl
var templates embed.FS

// Execute runs built-in and user-defined templates and links user defined files
func Execute(color models.Output) error {
	defer link()

	if err := os.MkdirAll(pathutil.StateDir, 0755); err != nil {
		slog.Error("Failed to create app cache directory", "error", err)
		return err
	}

	defaultTmpl := template.New("").Funcs(funcs)

	builtInTmpls, err := defaultTmpl.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return err
	}

	for _, tmpl := range builtInTmpls.Templates() {
		execute(tmpl, color)
	}

	templateRoot := filepath.Join(pathutil.ConfigDir, "templates")

	templatePath, err := filepath.Glob(templateRoot + "/*.tmpl")
	if err != nil {
		slog.Error("Failed to find template", "error", err)
		return err
	}

	if len(templatePath) == 0 {
		slog.Info("No user defined templates")
		return nil
	}

	userTmpl := template.New("").Funcs(funcs)
	userTmpls, err := userTmpl.ParseFiles(templatePath...)
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return err
	}

	for _, tmpl := range userTmpls.Templates() {
		execute(tmpl, color)
	}

	return nil
}

func link() error {
	links, err := config.GetStringMapStringSliceE("links")
	if err != nil {
		return fmt.Errorf("failed to convert links: %v", links)
	}

	for src, target := range links {
		for path := range slices.Values(target) {
			path, err := pathutil.FindPath(pathutil.ConfigDir, path)
			if err != nil {
				slog.Error("Failed to find path", "src", src, "path", path, "error", err)
				continue // log and continue instead of returning
			}
			if err := hardlinkOrCopy(filepath.Join(pathutil.StateDir, src), path); err != nil {
				slog.Error("Failed to link or copy", "src", src, "path", path, "error", err)
				continue
			}
		}
	}

	return nil
}

// execute executes a template using color
func execute(tmpl *template.Template, color models.Output) {
	name := tmpl.Name()
	saveFile := strings.TrimSuffix(name, ".tmpl")
	outputPath := filepath.Join(pathutil.StateDir, saveFile)

	if _, ok := success[saveFile]; ok {
		slog.Warn("Overwriting templates", "name", name)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		slog.Error("Error creating file", "file", outputPath, "error", err)
		success[saveFile] = false
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, color)
	if err != nil {
		slog.Error("Error executing template", "template", name, "error", err)
		success[saveFile] = false
		return
	}

	success[saveFile] = true
	slog.Info("Template written", "template", name, "path", outputPath)
}
