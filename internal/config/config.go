package config

import (
	"encoding/json"
	"os"
)

func Init(filename string) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&Conf)
	if err != nil {
		panic(err)
	}
}
