package config

import (
	"fmt"
	"os"
)

func GetApiPort() string {
	return fmt.Sprintf(":%s", os.Getenv("API_PORT"))
}
