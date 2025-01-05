package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// Находим корень проекта, поднимаясь по дереву директорий
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Проверяем наличие файла go.mod
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Если достигли корня файловой системы, прекращаем поиск
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", errors.New("не удалось найти корневую директорию проекта")
}
