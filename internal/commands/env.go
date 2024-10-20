package commands

import (
	"os"

	"github.com/joho/godotenv"
)

func getEnv(envFileName string, recursionLevels int) (map[string]string, error) {
	var mainEnv map[string]string
	var err error

	if envFileName == "none" {
		return mainEnv, nil
	}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Chdir(path)
	}()

	for i := 0; i < recursionLevels; i++ {
		mainEnv, err = godotenv.Read(envFileName)
		if err == nil {
			break
		}
		if err := os.Chdir(".."); err != nil {
			return nil, err
		}

		continue
	}

	return mainEnv, err
}
