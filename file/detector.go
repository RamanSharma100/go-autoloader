package file

import (
	"path/filepath"

	"github.com/RamanSharma100/go-autoloader/core"
)

var codeExtensions = map[string]bool{
	".go":   true,
	".js":   true,
	".ts":   true,
	".py":   true,
	".java": true,
	".rs":   true,
	".c":    true,
	".cpp":  true,
	".lua":  true,
}

var textExtensions = map[string]bool{
	".txt":  true,
	".md":   true,
	".env":  true,
	".yaml": true,
	".yml":  true,
	".toml": true,
}

var jsonExtensions = map[string]bool{
	".json": true,
}

var binaryExtensions = map[string]bool{
	".dll":  true,
	".so":   true,
	".exe":  true,
	".wasm": true,
	".png":  true,
	".jpg":  true,
}

func Detect(path string) core.FileKind {
	ext := filepath.Ext(path)

	switch {
	case codeExtensions[ext]:
		return core.CodeKind
	case textExtensions[ext]:
		return core.TextKind
	case jsonExtensions[ext]:
		return core.JSONKind
	case binaryExtensions[ext]:
		return core.BinaryKind
	default:
		return core.UnkownKind
	}
}
