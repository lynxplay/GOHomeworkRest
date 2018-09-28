package main

import (
	"encoding/json"
	"io/ioutil"
)

//RestServerConfiguration was outsourced
type RestServerConfiguration struct {
	SessionKey string

	DatabaseHost     string
	DatabasePort     int
	DatabaseUsername string
	DatabasePassword string
	Database         string

	DatabaseFileName string
}

//LoadServerConfiguration was outsourced
func LoadServerConfiguration(path string, defaultConfiguration RestServerConfiguration) (*RestServerConfiguration, error) {
	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		data, marshallError := json.MarshalIndent(defaultConfiguration, "", "	")
		if marshallError != nil {
			panic(marshallError)
		}

		ioutil.WriteFile(path, data, 0644)
		return &defaultConfiguration, nil
	}

	json.Unmarshal(bytes, &defaultConfiguration)
	return &defaultConfiguration, nil
}
