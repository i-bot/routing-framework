package settings

import (
	"encoding/json"
	"errorHandler"
	"os"
)

type Settings struct {
	Username, Password, IP, Port, DB_Name string

	ConnectionTaskStackSleepTime int

	Databases                                                                              [][]string
	Connections, ActionStack, OnOpenActions, OnWriteActions, OnReadActions, OnCloseActions string
}

var location = os.Getenv("HOME") + "/.network-distributor/settings.json"

func LoadSettings() (settings *Settings) {
	file, err := os.Open(location)
	errorHandler.HandleError(err)
	defer file.Close()

	decoder := json.NewDecoder(file)

	properties := Settings{}
	errorHandler.HandleError(decoder.Decode(&properties))

	return &properties
}
