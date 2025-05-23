package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Nadim147c/rong/internal/models"
)

//go:embed built-in/*.tmpl
var templates embed.FS

func Execute(output models.Output) {
	tmpls, err := template.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse templates: %v\n", err)
		return
	}

	// Get user cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user cache directory: %v\n", err)
		return
	}

	appCacheDir := filepath.Join(cacheDir, "rong")
	if err := os.MkdirAll(appCacheDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app cache directory: %v\n", err)
		return
	}

	for _, tmpl := range tmpls.Templates() {
		name := tmpl.Name()

		// Remove `.tmpl` suffix from filename
		saveFile := strings.TrimSuffix(name, ".tmpl")
		outputPath := filepath.Join(appCacheDir, saveFile)

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", outputPath, err)
			continue
		}

		err = tmpls.ExecuteTemplate(file, name, output)
		file.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing template %s: %v\n", name, err)
		} else {
			fmt.Printf("Template %s written to %s\n", name, outputPath)
		}
	}
}
