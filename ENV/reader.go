package ENV

import (
	"HostelApp/LogColor"
	"log/slog"
	"os"
	"strconv"
)

func ReadString(key string, value string) string {
	newValue := os.Getenv(key)
	if newValue != "" {
		return newValue
	}
	slog.Error(LogColor.Red("Error reading environment fail to get " + key + " from env file"))
	return value
}

func ReadInt64(key string, value int64) int64 {
	newValue := os.Getenv(key)
	if newValue != "" {
		parseInt, err := strconv.ParseInt(newValue, 10, 64)
		if err != nil {
			return value
		}
		return parseInt
	}
	slog.Error(LogColor.Red("Error reading environment fail to get " + key + " from env file"))
	return value
}

func ReadInt(key string, value int) int {
	newValue := os.Getenv(key)
	if newValue != "" {
		parseInt, err := strconv.Atoi(newValue)
		if err != nil {
			return value
		}
		return parseInt
	}
	slog.Error(LogColor.Red("Error reading environment fail to get " + key + " from env file"))
	return value
}
