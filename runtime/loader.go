package runtime

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"
)

type Loader struct {
	PluginDir string
	SourceDir string
}

func NewLoader(src, out string) *Loader {
	_ = os.MkdirAll(out, 0755)
	return &Loader{
		SourceDir: src,
		PluginDir: out,
	}
}

func (l *Loader) BuildAndLoad(filePath string) (string, error) {
	if runtime.GOOS == "windows" {
		return filePath, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, filePath, content, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse plugin structure: %w", err)
	}

	compilePath := filePath
	if fileAST.Name.Name != "main" {
		packageEndPos := fset.Position(fileAST.Name.End()).Offset
		remainingCode := string(content[packageEndPos:])
		rewrittenCode := "package main\n" + remainingCode

		tmpPluginFile := filepath.Join(l.PluginDir, "tmp_"+filepath.Base(filePath))
		if err := os.WriteFile(tmpPluginFile, []byte(rewrittenCode), 0644); err != nil {
			return "", fmt.Errorf("failed to create compile asset: %w", err)
		}
		compilePath = tmpPluginFile
		defer os.Remove(tmpPluginFile)
	}

	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(fileName, ".go")
	outPath := filepath.Join(l.PluginDir, name+".so")

	cmd := exec.Command(
		"go", "build",
		"-buildmode=plugin",
		"-o", outPath,
		compilePath,
	)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Go Build Failed: %v. Compiler Output: %s", err, stderr.String())
		return "", err
	}

	log.Println("[PLUGIN BUILT]", outPath)
	return outPath, nil
}

func (l *Loader) LoadPlugin(path string) (*plugin.Plugin, error) {
	if runtime.GOOS == "windows" {
		return nil, nil
	}
	return plugin.Open(path)
}

func (l *Loader) CallPlugin(p *plugin.Plugin, filePath string, name string, args ...any) error {
	if runtime.GOOS == "windows" || p == nil {
		output, err := l.CallWindowsStreamFallback(filePath, name, args...)
		if err != nil {
			return err
		}
		if output != "" {
			fmt.Print(output)
		}
		return nil
	}

	sym, err := p.Lookup(name)
	if err != nil {
		return err
	}

	fn, ok := sym.(func(...any))
	if !ok {
		return fmt.Errorf("symbol %s does not match expected function signature func(...any)", name)
	}

	fn(args...)
	return nil
}

func (l *Loader) CallWindowsStreamFallback(sourceFile string, targetName string, args ...any) (string, error) {
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, sourceFile, content, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse go file structure: %w", err)
	}

	packageEndPos := fset.Position(fileAST.Name.End()).Offset
	remainingCode := string(content[packageEndPos:])

	var execBlock string
	if len(args) > 0 {
		var argExprs []string
		for _, arg := range args {
			argExprs = append(argExprs, fmt.Sprintf("%#v", arg))
		}
		execBlock = fmt.Sprintf(`
	if fn, ok := interface{}(%s).(func(...interface{})); ok {
		fn([]interface{}{%s}...)
	} else if fn, ok := interface{}(%s).(func()); ok {
		fn()
	} else {
		fmt.Print(%s)
	}
`, targetName, strings.Join(argExprs, ", "), targetName, targetName)
	} else {
		execBlock = fmt.Sprintf(`
	if fn, ok := interface{}(%s).(func()); ok {
		fn()
	} else if fn, ok := interface{}(%s).(func(...interface{})); ok {
		fn()
	} else {
		fmt.Print(%s)
	}
`, targetName, targetName, targetName)
	}

	runnerTemplate := fmt.Sprintf(`package main

import "fmt"

%s

func main() {
	%s
}
`, remainingCode, execBlock)

	var stdout, stderr bytes.Buffer

	// Passing "-" as the file argument tells "go run" to read from stdin stream instead of disk
	cmd := exec.Command("go", "run", "-")

	// Stream the generated code template text completely out of your program's memory RAM pool
	cmd.Stdin = strings.NewReader(runnerTemplate)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execution failed: %s", stderr.String())
	}

	return stdout.String(), nil
}
