package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

type Configuration struct {
	Broker BrokerConfiguration
	Database DatabaseConfiguration
}

type BrokerConfiguration struct {
	Host string
	Port int
	ClientId string
	User string
	Pass string
}

type DatabaseConfiguration struct {
	ConnectionString string
	DevicesCollection string
	GroupsCollection string
	TemperatureReportsCollection string
	HumidityReportsCollection string
	PressureReportsCollection string
	IlluminanceReportsCollection string
}

var configurationFiles map[string]*Configuration

func GetConfig(params ...string) Configuration {

	configuration := Configuration{}
	env := "development"

	if len(params) > 0 {
		env = params[0]
	}

	_, filepath, _, ok := runtime.Caller(1)

	filename := fmt.Sprintf("../config/config.%s.json", env)

	if configurationFiles[filename] != nil {
		return *configurationFiles[filename]
	}

	absoluteFilename := path.Join(path.Dir(filepath), filename)

	if !ok {
		panic(fmt.Sprintf("Configuration file does not exist: %s", absoluteFilename))
	}

	configurationFile, err := os.Open(absoluteFilename)

	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configurationFile)
	jsonParser.Decode(&configuration)

	configurationFile.Close()

	cacheConfiguration(filename, configuration)

	return configuration
}

func cacheConfiguration(filename string, configuration Configuration) {

	if configurationFiles == nil {
		configurationFiles = make(map[string]*Configuration)
	}

	configurationFiles[filename] = &configuration

}