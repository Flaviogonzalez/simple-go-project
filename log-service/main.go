package main

import "log-service/config"

func main() {
	config.InitConfig().InitServer()
}
