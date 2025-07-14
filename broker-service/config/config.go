package config

import (
	"broker-service/routes"
	"log"
	"net/http"
)

type Config struct {
}

func InitConfig() *Config {

	return &Config{}
}

func (app *Config) StartServer() {

	server := &http.Server{
		Addr:    ":80",
		Handler: routes.Routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
