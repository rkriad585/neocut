//go:build ignore

package main

import (
	"os"
	"path/filepath"
	"strings"
)

func main() {
	src := filepath.Join("..", "..", ".version")
	dst := "version.txt"

	data, err := os.ReadFile(src)
	if err != nil {
		os.Exit(1)
	}
	if err := os.WriteFile(dst, []byte(strings.TrimSpace(string(data))+"\n"), 0644); err != nil {
		os.Exit(1)
	}
}
