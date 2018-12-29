package main

import (
	app "github.com/j33ty/iota-service/includes/app"
	utils "github.com/j33ty/iota-service/includes/utils"
	services "github.com/j33ty/iota-service/services"
)

func main() {

	appName := "iota-service"
	c := app.GetInstance()
	err := c.Initialize(appName)
	utils.Err("Initialize app", err)

	services.Run(c)

	c.Shutdown()
}
