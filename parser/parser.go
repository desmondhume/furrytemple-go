package parser

import (
	"encoding/json"
	"fmt"
)

func parseYoutube(data []byte) map[string]interface{} {
	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	return parsed
}

func Parse(source string, input chan []byte, output chan map[string]interface{}) {
	var parseFunc func([]byte) map[string]interface{}

	switch source {
	case "youtube":
		parseFunc = parseYoutube
	default:
		fmt.Println("Adatper not found")
	}

	for data := range input {
		output <- parseFunc(data)
	}
}
