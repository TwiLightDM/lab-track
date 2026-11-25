package main

import (
	"github.com/TwiLightDM/lab-track/internal/app"
	"github.com/TwiLightDM/lab-track/internal/config"
)

// @title        lab track api
// @version      1.0
// @description  API для управления лабораторным оборудованием
// @host         localhost:8080/api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}
func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
