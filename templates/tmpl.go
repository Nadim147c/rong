package templates

import (
	"embed"
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
func Execute(color models.Output) {
	defer link()

	if err := os.MkdirAll(pathutil.StateDir, 0755); err != nil {
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

	templateRoot := filepath.Join(pathutil.ConfigDir, "templates")

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
	links := viper.Get("links")
	if links == nil {
		slog.Warn("No links defined in configuration")
		return
	}

	linksValue := reflect.ValueOf(links)
	if linksValue.Kind() != reflect.Map {
		slog.Error("Expected 'links' to be a map in configuration")
		return
	}

	for _, key := range linksValue.MapKeys() {
		if key.Kind() != reflect.String {
			slog.Error("Link key is not a string", "key", key)
			os.Exit(1) // Fatal configuration error
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
				os.Exit(1)
			}
			srcDir := filepath.Join(pathutil.StateDir, src)
			hardlinkOrCopy(srcDir, path)

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
					os.Exit(1)
				}
				srcDir := filepath.Join(pathutil.StateDir, src)
				hardlinkOrCopy(srcDir, path)
			}

		default:
			slog.Error("Invalid path value type; expected string or array/slice", "src", src, "type", pathValue.Kind().String())
			os.Exit(1)
		}
	}
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
