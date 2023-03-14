package config

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	Http struct {
		Addr          string `json:"port"`
		ReadTimeout   int    `json:"readTimeout"`
		WriteTimeout  int    `json:"writeTimeout"`
		MaxHeaderByte int    `json:"oneMB"`
	}

	Db struct {
		Driver     string `json:"driver"`
		DBName     string `json:"dbname"`
		CtxTimeout int    `json:"ctxTimeout"`
	}
}

func NewConfig(cfgFilePath string) (*Config, error) {
	var cfg *Config
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return &Config{}, err
	}
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&cfg); err != nil {
		return &Config{}, err
	}
	return cfg, nil
}
