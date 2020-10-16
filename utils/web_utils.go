package utils

import "os"

type DecodeResult struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

func DecodeResultJson(status int, message string) DecodeResult {
	return DecodeResult{
		Status: status,
		Text:   message,
	}
}

func GetConfigFile() string {
	currentPath, _ := os.Getwd()
	return currentPath + "/config.yaml"
}
