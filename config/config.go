// Package config provides functions for reading and writing JSON config files.
package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// ReadConfig reads a JSON config file and unmarshals it into the provided struct.
func ReadConfig(fileName string, config interface{}) error {
	configFile, err := os.Open(fileName)
	if err != nil {
		return errors.New("error opening file")
	}
	defer configFile.Close()

	fileData, err := io.ReadAll(configFile)
	if err != nil {
		return errors.New("error reading file")
	}

	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return errors.New("error unmarshalling file")
	}

	return nil
}
