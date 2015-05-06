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
	Connections, ActionQueue, OnOpenActions, OnWriteActions, OnReadActions, OnCloseActions string
}

func LoadSettings(location string) (settings *Settings) {
	file, err := os.Open(location)
	errorHandler.HandleError(err)
	defer file.Close()

	decoder := json.NewDecoder(file)

	properties := Settings{}
	errorHandler.HandleError(decoder.Decode(&properties))

	return &properties
}
