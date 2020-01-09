package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server
}

var Settings *Config = &Config{}

type Server struct {
	GrpcHost     string
	GrpcPort     int32
	GatewayHost  string
	GatewayPort  int32
}

func InitConfig(path string, config interface{}) {
	configPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error happened: %v \n", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Printf("Error happened: %v \n", err)
		os.Exit(1)
	}
}


func Insert(slice, insertion []interface{}, index int) []interface{} {
	result := make([]interface{}, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

