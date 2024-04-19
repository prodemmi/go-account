package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func GetTestBaseUrl(path string) string {
	return fmt.Sprintf("http://%s:%s/v1/%s/", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"), path)
}

func GetTestPingUrl() string {
	return fmt.Sprintf("http://%s:%s/ping", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
}

func StructToJsonBody(s interface{}) *bytes.Buffer {
	jsonValue, _ := json.Marshal(s)
	return bytes.NewBuffer(jsonValue)
}
