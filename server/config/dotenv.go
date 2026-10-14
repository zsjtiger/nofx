package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadDotEnv loads environment variables from the first .env file it finds in
// the working directory or its parent. The server is usually started from
// server/ while .env lives at the repository root, so both places are checked.
// A missing file is not an error.
func LoadDotEnv() error {
	_, err := loadFirstDotEnv(".env", filepath.Join("..", ".env"))
	return err
}

// loadFirstDotEnv loads the first existing file among candidates and returns
// its path. It returns an empty path when none exists.
func loadFirstDotEnv(candidates ...string) (string, error) {
	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}
		return path, godotenv.Load(path)
	}
	return "", nil
}
