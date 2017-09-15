package main

import (
	"io/ioutil"
	"encoding/json"
)

type Load struct {
	Regexp     string
	GroupNames []string
	TimestampFormat string
}

type Filter struct {
	FieldName string
	Regexp    string
}

type Config struct {
	Load    []Load
	Filters []Filter
	OutputFields []string
	MiscOptions MiscOptions
}

type MiscOptions struct {
	SpaceReplacement string
	OmitIfEmpty bool
}

func loadConfig(filename string) (Config, error){
	var config Config
	configData, err := ioutil.ReadFile(filename)

	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(configData, &config); err != nil {
		return config, err
	}

	return config, nil
}
