package config

import (
	"fmt"
	"os"
)

func GetApplicationPort() string {
	return fmt.Sprintf(":%s", os.Getenv("API_PORT"))
}

func GetRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}

func GetRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func GetDynamoDbTableName() string {
	return os.Getenv("TABLE_NAME")
}
