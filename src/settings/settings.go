package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

//var settingsHolder SettingsHolder

type Settings struct {
	Username, Password, IP, Port, DB_Name string
}

var location = os.Getenv("HOME") + "/.network-distributor/settings.json"

func LoadSettings() (settings Settings) {
	file, err := os.Open(location)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	settings = Settings{}
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return
	}
	
	return
}
