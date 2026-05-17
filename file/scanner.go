package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/RamanSharma100/go-autoloader/core"
)

func Scan(root string) ([]*core.Module, error) {
	var modules []*core.Module

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}

	err = filepath.Walk(rootAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if seen[path] {
			return nil
		}
		seen[path] = true

		content, _ := os.ReadFile(path)

		fileName := filepath.Base(path)
		ext := filepath.Ext(path)
		name := strings.TrimSuffix(fileName, ext)

		rel, _ := filepath.Rel(rootAbs, path)
		dir := filepath.Dir(rel)

		dir = strings.ReplaceAll(dir, "\\", ".")
		dir = strings.ReplaceAll(dir, "/", ".")
		dir = strings.Trim(dir, ".")

		runtimePath := name
		if dir != "." && dir != "" {
			runtimePath = dir + "." + name
		}

		module := &core.Module{
			Name:        name,
			Path:        path,
			Extension:   ext,
			Kind:        Detect(path),
			Content:     content,
			RuntimeName: name,
			RuntimePath: runtimePath,
			RelativeDir: dir,
		}

		modules = append(modules, module)
		return nil
	})

	return modules, err
}
