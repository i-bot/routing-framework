package settings

import (
	"encoding/json"
	"errorHandler"
	"os"
)

type Settings struct {
	Username, Password, IP, Port, DB_Name string

	Databases [][]string
}

var location = os.Getenv("HOME") + "/.network-distributor/settings.json"

func LoadSettings() (settings Settings) {
	file, err := os.Open(location)
	errorHandler.HandleError(err)
	defer file.Close()

	decoder := json.NewDecoder(file)

	settings = Settings{}
	errorHandler.HandleError(decoder.Decode(&settings))

	return
}
