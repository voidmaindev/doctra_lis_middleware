package websockets

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func HandleWSUpgradeMiddleware(c *fiber.Ctx) error {
	fmt.Println("Upgrading to websocket")

	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}	

	fmt.Println("Cant upgrade to websocket")

	return fiber.ErrUpgradeRequired
}

func HandleWS(c *websocket.Conn) {
	fmt.Println("Handling websocket")

	var (
		mt int
		message []byte
		err error
	)

	for {
		mt, message, err = c.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		err = c.Conn.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}