package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Abirdcfly/md-mermaid-silent-check-fix/model"
)

func ScanDirectory(root string) ([]model.MarkdownFile, error) {
	var files []model.MarkdownFile

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		files = append(files, model.MarkdownFile{
			Path:    path,
			Content: string(content),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
