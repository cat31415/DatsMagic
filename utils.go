package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func saveStateToJSON(state *MoveResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // для удобства чтения JSON

	if err := encoder.Encode(state); err != nil {
		return fmt.Errorf("ошибка при кодировании JSON: %v", err)
	}

	fmt.Println("Данные успешно сохранены в", filename)
	return nil
}
