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

	"github.com/Nadim147c/rong/v3/internal/models"
	"github.com/Nadim147c/rong/v3/internal/pathutil"
	"github.com/google/renameio/v2"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var success = map[string]bool{}

//go:embed built-in/*.tmpl
var templates embed.FS

// Execute runs built-in and user-defined templates and links user defined files
func Execute(color models.Output) error {
	defer link()

	if err := os.MkdirAll(pathutil.StateDir, 0o755); err != nil {
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
	rawLinks := viper.Get("links")
	links, err := cast.ToStringMapStringSliceE(rawLinks)
	if err != nil {
		return fmt.Errorf("failed to convert links: %v", links)
	}

	for name, target := range links {
		for path := range slices.Values(target) {
			dst, err := pathutil.FindPath(pathutil.ConfigDir, path)
			if err != nil {
				slog.Error(
					"Failed to find path",
					"src", name,
					"path", dst,
					"error", err,
				)
				continue // log and continue instead of returning
			}
			src := filepath.Join(pathutil.StateDir, name)
			if err := atomicCopy(src, dst); err != nil {
				slog.Error(
					"Failed to atomically copy",
					"src", name,
					"path", dst,
					"error", err,
				)
				continue
			}
			slog.Info("Successfully copied", "src", name, "path", dst)
		}
	}

	return nil
}

// execute executes a template using color
func execute(tmpl *template.Template, out models.Output) {
	name := tmpl.Name()
	saveFile := strings.TrimSuffix(name, ".tmpl")
	outputPath := filepath.Join(pathutil.StateDir, saveFile)

	if _, ok := success[saveFile]; ok {
		slog.Warn("Overwriting templates", "name", name)
	}

	f, err := renameio.TempFile("", outputPath)
	if err != nil {
		slog.Error("Error executing template", "template", name, "error", err)
		success[saveFile] = false
		return
	}
	defer f.Cleanup()

	if err := tmpl.Execute(f, out); err != nil {
		slog.Error("Error executing template", "template", name, "error", err)
		success[saveFile] = false
		return
	}

	if err := f.CloseAtomicallyReplace(); err != nil {
		slog.Error(
			"Failed atomic of replace",
			"file", outputPath,
			"error", err,
		)
		success[saveFile] = false
		return
	}

	success[saveFile] = true
	slog.Info("Template written", "template", name, "path", outputPath)
}
