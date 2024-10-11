package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	// Отправка GET запроса
	resp, err := http.Get("http://example.com/data.json")
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Декодирование JSON
	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Ошибка при декодировании JSON:", err)
		return
	}

	// Использование данных
	fmt.Println("Ключ:", data.Key)
	fmt.Println("Значение:", data.Value)
}
