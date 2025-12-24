package templates

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"text/template"
	"unicode"

	"github.com/MatusOllah/stripansi"
	"github.com/Nadim147c/rong/v4/internal/models"
	"github.com/Nadim147c/rong/v4/internal/pathutil"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var success = successCounter{}

type successCounter map[string]struct{}

func (s successCounter) set(n string) { s[n] = struct{}{} }
func (s successCounter) has(n string) bool {
	_, ok := s[n]
	return ok
}

//go:embed built-in/*.tmpl
var templates embed.FS

// Execute runs built-in and user-defined templates and links user defined files
func Execute(colors models.Output) error {
	var allErrors []error

	if err := os.MkdirAll(pathutil.StateDir, 0o755); err != nil {
		slog.Error("Failed to create app cache directory", "error", err)
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	defaultTmpl := template.New("").Funcs(funcs)

	builtInTmpls, err := defaultTmpl.ParseFS(templates, "built-in/*.tmpl")
	if err != nil {
		slog.Error("Failed to parse templates", "error", err)
		return fmt.Errorf("failed to parse built-in templates: %w", err)
	}

	// Execute built-in templates and collect errors
	for _, tmpl := range builtInTmpls.Templates() {
		if err := execute(tmpl, colors); err != nil {
			allErrors = append(allErrors, err)
			slog.Error(
				"Error executing built-in template",
				"template",
				tmpl.Name(),
				"error",
				err,
			)
		}
	}

	templateRoot := filepath.Join(pathutil.ConfigDir, "templates")
	templatePaths, err := filepath.Glob(templateRoot + "/*.tmpl")
	if err != nil {
		slog.Error("Failed to find templates", "error", err)
		allErrors = append(
			allErrors,
			fmt.Errorf("failed to find user templates: %w", err),
		)
	} else if len(templatePaths) == 0 {
		slog.Info("No user defined templates")
	} else {
		userTmpl := template.New("").Funcs(funcs)
		userTmpls, err := userTmpl.ParseFiles(templatePaths...)
		if err != nil {
			slog.Error("Failed to parse user templates", "error", err)
			allErrors = append(allErrors, fmt.Errorf("failed to parse user templates: %w", err))
		} else {
			// Execute user templates and collect errors
			for _, tmpl := range userTmpls.Templates() {
				if err := execute(tmpl, colors); err != nil {
					allErrors = append(allErrors, err)
					slog.Error("Error executing user template", "template", tmpl.Name(), "error", err)
				}
			}
		}
	}

	// Run post-hook and collect any errors
	postHookErrs := postHook(colors)
	if postHookErrs != nil {
		allErrors = append(allErrors, postHookErrs)
	}

	// Return combined errors if any occurred
	if len(allErrors) > 0 {
		return fmt.Errorf(
			"template execution completed with errors: %w",
			errors.Join(allErrors...),
		)
	}

	return nil
}

func getConfig(key string) (map[string][]string, error) {
	rawCfg := viper.Get(key)
	if rawCfg == nil {
		return map[string][]string{}, nil // user has not specified config
	}
	return cast.ToStringMapStringSliceE(rawCfg)
}

func postHook(colors models.Output) error {
	var allErrors []error

	linksCfg, err := getConfig("links")
	if err != nil {
		allErrors = append(
			allErrors,
			fmt.Errorf("failed to parse links config: %w", err),
		)
	}

	installsCfg, err := getConfig("installs")
	if err != nil {
		allErrors = append(
			allErrors,
			fmt.Errorf("failed to parse installs config: %w", err),
		)
	}

	cmdsCfg, err := getConfig("post-cmds")
	if err != nil {
		allErrors = append(
			allErrors,
			fmt.Errorf("failed to parse post_hooks config: %w", err),
		)
	}

	exe, err := os.Executable()
	if err != nil {
		allErrors = append(
			allErrors,
			fmt.Errorf("failed to get executable path: %w", err),
		)
	}

	// Base environment variables
	baseEnv := os.Environ()
	baseEnviron := map[string]string{
		"RONG":        exe,
		"RONG_CACHE":  pathutil.CacheDir,
		"RONG_CONFIG": pathutil.ConfigDir,
		"RONG_STATE":  pathutil.StateDir,
		"IMAGE":       colors.Image,
		"RONG_DARK":   viper.GetString("dark"),
		// material envs
		"RONG_MATERIAL_VARIANT":  viper.GetString("material.variant"),
		"RONG_MATERIAL_VERSION":  viper.GetString("material.version"),
		"RONG_MATERIAL_CONTRAST": viper.GetString("material.contrast"),
		"RONG_MATERIAL_PLATFORM": viper.GetString("material.platform"),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	addErrs := func(errs ...error) {
		mu.Lock()
		allErrors = append(allErrors, errs...)
		mu.Unlock()
	}

	lock := make(chan struct{}, runtime.NumCPU())
	for name := range success {
		lock <- struct{}{} // avoid creating too many goroutines
		wg.Go(func() {
			defer func() { <-lock }() // release the lock

			// Build environment for this specific template
			cmdEnv := make([]string, 0, len(baseEnv)+len(baseEnviron)+5)
			cmdEnv = append(cmdEnv, baseEnv...)

			// Add base environment variables
			for k, v := range baseEnviron {
				cmdEnv = addEnv(cmdEnv, k, v)
			}

			// Set source file information
			sourcePath := filepath.Join(pathutil.StateDir, name)
			cmdEnv = addEnv(cmdEnv, "RONG_SOURCE", sourcePath)
			cmdEnv = addEnv(cmdEnv, "RONG_SOURCE_NAME", name)

			// Track installed and linked paths
			var installedPaths, linkedPaths []string

			// Process links
			if links, ok := linksCfg[name]; ok && links != nil {
				linked, err := link(name, links)
				if err != nil {
					addErrs(fmt.Errorf("failed to link %s: %w", name, err))
				}
				linkedPaths = linked
			}

			// Process installs
			if installs, ok := installsCfg[name]; ok && installs != nil {
				installed, err := install(name, installs)
				if err != nil {
					addErrs(fmt.Errorf("failed to install %s: %w", name, err))
				}
				installedPaths = installed
			}

			cmdEnv = addEnvPaths(cmdEnv, "RONG_INSTALLED", installedPaths)
			cmdEnv = addEnvPaths(cmdEnv, "RONG_LINKED", linkedPaths)
			copied := slices.Concat(installedPaths, linkedPaths)
			cmdEnv = addEnvPaths(cmdEnv, "RONG_COPIED", copied)

			// Run hooks with the complete environment
			if hooks, ok := cmdsCfg[name]; ok && hooks != nil {
				if err := runHooks(name, hooks, cmdEnv); err != nil {
					addErrs(
						fmt.Errorf("failed to run hooks for %s: %w", name, err),
					)
				}
			}
		})
	}

	wg.Wait()

	return errors.Join(allErrors...)
}

func runHooks(name string, hooks []string, env []string) error {
	var errs []error

	for hook := range slices.Values(hooks) {
		cmd := exec.Command("sh", "-c", hook)
		cmd.Env = env

		hook = strings.TrimRightFunc(hook, unicode.IsSpace)

		out, err := cmd.CombinedOutput()
		out = stripansi.Bytes(out)
		out = bytes.TrimRightFunc(out, unicode.IsSpace)
		if err != nil {
			slog.Error(
				"Hook failed",
				"name", name,
				"hook", hook,
				"output", string(out),
				"error", err,
			)
			errs = append(
				errs,
				fmt.Errorf("hook %q failed: %w (output: %s)", hook, err, out),
			)
			continue
		}

		slog.Info(
			"Hook executed successfully",
			"name", name,
			"hook", hook,
			"output", string(out),
		)
	}

	return errors.Join(errs...)
}

func addEnv(envs []string, key, val string) []string {
	return append(envs, fmt.Sprintf("%s=%s", key, val))
}

func addEnvPaths(envs []string, key string, vals []string) []string {
	return append(envs, fmt.Sprintf("%s=%s", key, strings.Join(vals, ":")))
}

func install(src string, targets []string) ([]string, error) {
	var errs []error
	var installedPaths []string

	for _, path := range targets {
		dst, err := pathutil.FindPath(pathutil.ConfigDir, path)
		if err != nil {
			slog.Error(
				"Failed to find path",
				"src", src,
				"path", path,
				"error", err,
			)
			errs = append(
				errs,
				fmt.Errorf("failed to find path %q: %w", path, err),
			)
			continue
		}
		srcPath := filepath.Join(pathutil.StateDir, src)
		if err := atomicCopy(srcPath, dst); err != nil {
			slog.Error(
				"Failed to atomically copy",
				"src", srcPath,
				"dst", dst,
				"error", err,
			)
			errs = append(
				errs,
				fmt.Errorf("failed to copy %q to %q: %w", srcPath, dst, err),
			)
			continue
		}
		installedPaths = append(installedPaths, dst)
		slog.Info("Successfully copied", "src", srcPath, "dst", dst)
	}

	if len(errs) > 0 {
		return installedPaths, errors.Join(errs...)
	}
	return installedPaths, nil
}

func link(src string, targets []string) ([]string, error) {
	var errs []error
	var linkedPaths []string

	for _, path := range targets {
		dst, err := pathutil.FindPath(pathutil.ConfigDir, path)
		if err != nil {
			slog.Error(
				"Failed to find path",
				"src", src,
				"path", path,
				"error", err,
			)
			errs = append(
				errs,
				fmt.Errorf("failed to find path %q: %w", path, err),
			)
			continue
		}
		srcPath := filepath.Join(pathutil.StateDir, src)
		if err := hardlinkOrCopy(srcPath, dst); err != nil {
			slog.Error(
				"Failed to create link or copy",
				"src", srcPath,
				"dst", dst,
				"error", err,
			)
			errs = append(
				errs,
				fmt.Errorf(
					"failed to link/copy %q to %q: %w",
					srcPath, dst, err,
				),
			)
			continue
		}
		linkedPaths = append(linkedPaths, dst)
		slog.Info("Successfully linked/copied", "src", srcPath, "dst", dst)
	}

	if len(errs) > 0 {
		return linkedPaths, errors.Join(errs...)
	}
	return linkedPaths, nil
}

// execute executes a template using color and returns any error
func execute(tmpl *template.Template, out models.Output) error {
	name := tmpl.Name()
	filename := strings.TrimSuffix(name, ".tmpl")
	outputPath := filepath.Join(pathutil.StateDir, filename)

	if success.has(filename) {
		slog.Warn("Overwriting template", "name", name)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(
			"failed to create output file for template %q: %w",
			name, err,
		)
	}
	defer file.Close()

	if err := tmpl.Execute(file, out); err != nil {
		return fmt.Errorf("failed to execute template %q: %w", name, err)
	}

	success.set(filename)
	slog.Info("Template written", "template", name, "path", outputPath)
	return nil
}
