package utils

import (
	"encoding/json"
	"log"
)

func GetAsJsonString(obj any) string {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("error occurred while forming json string", err)
	}
	return string(data)
}
