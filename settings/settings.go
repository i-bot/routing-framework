package settings

import (
	"encoding/json"
	"errorHandler"
	"os"
)

type Settings struct {
	Username, Password, IP, Port, DB_Name string

	ActionQueueSleepTime int

	Connections, ActionQueue, OnOpen, OnWrite, OnRead, OnClose string
	Databases, Values                                          [][]string
}

func LoadSettings(location string) (properties *Settings) {
	file, err := os.Open(location)
	errorHandler.HandleError(err)
	defer file.Close()

	decoder := json.NewDecoder(file)

	properties = &Settings{}
	errorHandler.HandleError(decoder.Decode(&properties))

	return
}

func SaveSettings(location string, properties *Settings) {
	file, err := os.Create(location)
	errorHandler.HandleError(err)
	defer file.Close()

	encoder := json.NewEncoder(file)

	errorHandler.HandleError(encoder.Encode(&properties))
}
