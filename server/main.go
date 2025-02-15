package main

import (
	"github.com/bluefalconhd/lbd_game/server/config"
	"github.com/bluefalconhd/lbd_game/server/database"
	"github.com/bluefalconhd/lbd_game/server/routes"
	"github.com/bluefalconhd/lbd_game/server/utils"

	"github.com/joho/godotenv"
)

func main() {
	database.ConnectDatabase()
	utils.InitScheduler()

	godotenv.Load()
	cfg := config.LoadConfig()

	router := routes.SetupRouter(cfg)
	router.Run(":8040")
}
