package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of your YAML data
type Config struct {
	Port        int       `yaml:"port"`
	SensorID    string    `yaml:"SensorID"`
	OA_Axial    []float32 `yaml:"OA_Axial"`
	OA_Radial_1 []float32 `yaml:"OA_Radial_1"`
	OA_Radial_2 []float32 `yaml:"OA_Radial_2"`
	Def_Bea     []float32 `yaml:"Def_Bea"`
	Def_Imb     []float32 `yaml:"Def_Imb"`
	OA_QC       []float32 `yaml:"OA_QC"`
}

// writeConfig writes the Config struct to a YAML file
func writeConfig(filename string, config Config) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// readConfig reads the YAML file into a Config struct
func readConfig(filename string) (Config, error) {
	var config Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("error reading file: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	return config, nil
}
