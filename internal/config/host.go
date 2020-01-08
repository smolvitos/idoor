package config

import (
	"os"
)

func GetHost() string {
	env := os.Getenv("HOST")
	if len(env) == 0 {
		env = ":8080"
	}
	return env
}
