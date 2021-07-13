package util

import (
	"path/filepath"
)

func RelativePath(basePath, fullPath string) string {
	rel, err := filepath.Rel(basePath, fullPath)
	if err != nil {
		return fullPath
	}

	rel = filepath.ToSlash(rel)
	return rel
}
