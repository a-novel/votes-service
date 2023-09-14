package config

import (
	_ "embed"
	"log"
)

//go:embed api-dev.yml
var apiDevFile []byte

//go:embed api-prod.yml
var apiProdFile []byte

type ApiConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var API *ApiConfig

func init() {
	cfg := new(ApiConfig)

	if err := loadEnv(EnvLoader{ProdENV: apiProdFile, DevENV: apiDevFile}, cfg); err != nil {
		log.Fatalf("error loading api configuration: %v\n", err)
	}

	API = cfg
}
