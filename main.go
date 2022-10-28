package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tinkerbaj/chatwebsocketgin/chat"
	"github.com/tinkerbaj/chatwebsocketgin/middleware"
)

func main() {

	//create new Hub and run it
	hub := chat.NewHub()
	go hub.Run()
	go func() { // infinite loop to send notification
		for range time.Tick(time.Second * 30) {
			go hub.SendNotification("room1", "test")
		}
	}()
	//we need pass hub to out route with roomid
	app := gin.Default()
	app.Use(middleware.Next())
	app.Use(middleware.Count())
	app.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		chat.ServeWS(c, roomId, hub)
	})
	// _ = app.Run(":20192")
	app.Run()
}
