package web_fs_repository

import (
	"fmt"
	"os"

	core_error "github.com/FeeRro/golang-todoApp/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file: %s: %w", filePath, core_error.ErrNotFound)
		}
		return nil, fmt.Errorf("get file: %s, %w", filePath, err)
	}

	return file, nil
}
