package templates

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/Nadim147c/rong/internal/models"
	"github.com/Nadim147c/rong/internal/pathutil"
	"github.com/spf13/viper"
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

	tmpls, err := defaultTmpl.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return err
	}

	for _, tmpl := range tmpls.Templates() {
		execute(tmpl, color)
	}

	templateRoot := filepath.Join(pathutil.ConfigDir, "templates")

	templatePath, err := filepath.Glob(templateRoot + "/*.tmpl")
	if err != nil {
		slog.Error("Failed to find template", "error", err)
		return err
	}

	if len(templatePath) == 0 {
		return nil
	}

	tmpls, err = defaultTmpl.ParseFiles(templatePath...)
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return err
	}

	for _, tmpl := range tmpls.Templates() {
		execute(tmpl, color)
	}

	return nil
}

func link() error {
	links := viper.Get("links")
	if links == nil {
		slog.Warn("No links defined in configuration")
		return nil
	}

	linksValue := reflect.ValueOf(links)
	if linksValue.Kind() != reflect.Map {
		return errors.New("expected 'links' to be a map in configuration")
	}

	for _, key := range linksValue.MapKeys() {
		if key.Kind() != reflect.String {
			return fmt.Errorf("link key is not a string: %v", key)
		}

		src := key.String()
		if _, ok := success[src]; !ok {
			slog.Warn("Skipping source, it doesn't exist", "src", src)
			continue
		}
		if !success[src] {
			slog.Warn("Skipping source, it previously failed", "src", src)
			continue
		}

		pathValue := linksValue.MapIndex(key)
		if pathValue.Kind() == reflect.Interface {
			pathValue = pathValue.Elem()
		}

		switch pathValue.Kind() {
		case reflect.String:
			path, err := pathutil.FindPath(pathutil.ConfigDir, pathValue.String())
			if err != nil {
				slog.Error("Failed to find path", "src", src, "path", pathValue.String(), "error", err)
				continue
			}
			srcDir := filepath.Join(pathutil.StateDir, src)
			if err := hardlinkOrCopy(srcDir, path); err != nil {
				slog.Error("Failed to link or copy", "src", src, "path", path, "error", err)
				continue
			}

		case reflect.Array, reflect.Slice:
			for i := 0; i < pathValue.Len(); i++ {
				k := pathValue.Index(i)
				if k.Kind() == reflect.Interface {
					k = k.Elem()
				}

				if k.Kind() != reflect.String {
					slog.Warn("Skipping non-string path entry", "index", i, "type", k.Kind(), "src", src)
					continue
				}
				pathStr := k.String()
				path, err := pathutil.FindPath(pathutil.ConfigDir, pathStr)
				if err != nil {
					slog.Error("Failed to find path", "src", src, "path", pathStr, "error", err)
					continue // log and continue instead of returning
				}
				if err := hardlinkOrCopy(filepath.Join(pathutil.StateDir, src), path); err != nil {
					slog.Error("Failed to link or copy", "src", src, "path", path, "error", err)
					continue
				}
			}

		default:
			return fmt.Errorf("invalid path value type for src %q; expected string or array/slice, got %s",
				src, pathValue.Kind().String())
		}
	}

	return nil
}

// execute executes a template using color
func execute(tmpl *template.Template, color models.Output) {
	name := tmpl.Name()
	saveFile := strings.TrimSuffix(name, ".tmpl")
	outputPath := filepath.Join(pathutil.StateDir, saveFile)

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
