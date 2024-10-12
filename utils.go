package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func saveToJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Ошибка при кодировании JSON: %v", err)
	}

	if err := ioutil.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("Ошибка при записи в файл: %v", err)
	}

	return nil
}
