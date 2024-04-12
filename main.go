package main

import (
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/voidmaindev/doctra_lis_middleware/controllers"
	"github.com/voidmaindev/doctra_lis_middleware/inits"
	"github.com/voidmaindev/doctra_lis_middleware/store"
	"github.com/voidmaindev/doctra_lis_middleware/websockets"
)

var addr string

func init() {
	inits.LoadEnvVars()
	addr = os.Getenv("ADDRESS") + ":" + os.Getenv("PORT")

	store.ConnectToDB()
}

func main() {
	app := fiber.New()

	app.Get("/hardwares", controllers.HardwaresGetAll)
	app.Get("/hardwares/:id", controllers.HardwaresGetByID)
	app.Post("/hardwares", controllers.HardwaresCreate)
	app.Put("/hardwares/:id", controllers.HardwaresUpdate)
	app.Delete("/hardwares/:id", controllers.HardwaresDelete)

	app.Use("/", websockets.HandleWSUpgradeMiddleware)
	app.Get("/", websocket.New(websockets.HandleWS))

	app.Listen(addr)
}