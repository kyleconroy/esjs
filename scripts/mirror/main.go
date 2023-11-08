package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	// Assume it exists
	loc := flag.Arg(0)

	pkgs := []string{
		"ast",
		"compat",
		"config",
		"css_ast",
		"css_lexer",
		"helpers",
		"fs",
		"js_ast",
		"js_lexer",
		"js_parser",
		"js_printer",
		"logger",
		"renamer",
		"runtime",
		"sourcemap",
		"test",
	}

	dir := filepath.Join(loc, "internal")
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		var keep bool
		for _, pkg := range pkgs {
			if strings.Contains(path, filepath.Join("internal", pkg)) {
				keep = true
			}
		}

		if !keep {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newpath := strings.Replace(path, dir, ".", 1)

		os.MkdirAll(filepath.Dir(newpath), 0755)

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/evanw/esbuild/internal`),
			[]byte(`"github.com/kyleconroy/esjs`))

		if err := os.WriteFile(newpath, contents, 0644); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}

	for _, filename := range []string{"LICENSE.md", ".gitignore"} {
		path := filepath.Join(loc, filename)
		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(filename, contents, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
