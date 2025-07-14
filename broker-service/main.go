package main

import "broker-service/config"

func main() {
	config.InitConfig().StartServer()
}
