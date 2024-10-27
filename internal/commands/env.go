package commands

import (
	"os"

	"github.com/joho/godotenv"
)

func getEnv(envFileName string, path string, recursionLevels int) (map[string]string, error) {
	var mainEnv map[string]string
	var err error

	startPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Chdir(startPath)
	}()

	if path == "" {
		path = startPath
	}

	if err := os.Chdir(path); err != nil {
		return nil, err
	}

	for i := 0; i <= recursionLevels; i++ {
		mainEnv, err = godotenv.Read(envFileName)
		if err == nil {
			return mainEnv, nil
		}

		if err := os.Chdir(".."); err != nil {
			return nil, err
		}
	}

	return nil, ErrEnvNotFound
}
