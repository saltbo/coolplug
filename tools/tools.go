package tools

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkAllFile(dir string, suffix string, visit func(filename string) error) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			return visit(path)
		}

		return nil
	})
}
