package utils

import (
	"os"
)

func GetStoragePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./tasks.json"
	}
	return home + "/.neon/tasks.json"
}
