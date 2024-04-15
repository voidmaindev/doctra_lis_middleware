package main

import (
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/voidmaindev/doctra_lis_middleware/controllers"
	_ "github.com/voidmaindev/doctra_lis_middleware/docs"
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

// @title           Doctra Middleware API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:9000
// @BasePath  /
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/hardwares", controllers.HardwaresGetAll)
	app.Get("/hardwares/:id", controllers.HardwaresGetByID)
	app.Post("/hardwares", controllers.HardwaresCreate)
	app.Put("/hardwares/:id", controllers.HardwaresUpdate)
	app.Delete("/hardwares/:id", controllers.HardwaresDelete)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use("/", websockets.HandleWSUpgradeMiddleware)
	app.Get("/", websocket.New(websockets.HandleWS))

	app.Listen(addr)
}
