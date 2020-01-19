package config

import (
	"os"
)

func GetCode() string {
	env := os.Getenv("CODE")
	if len(env) == 0 {
		env = "Код: h2q4v2f"
	}
	return env
}
