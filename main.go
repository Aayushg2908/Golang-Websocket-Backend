package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var roomClients = make(map[string]map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func wsEndPoint(c *gin.Context) {
	fmt.Println("Connected")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	clients[ws] = true

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		for client := range clients {
			err := client.WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func roomMessage(c *gin.Context) {
	fmt.Println("Connected")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	roomName := c.Param("roomName")
	if roomClients[roomName] == nil {
		roomClients[roomName] = make(map[*websocket.Conn]bool)
	}

	roomClients[roomName][ws] = true

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		for client := range roomClients[roomName] {
			err := client.WriteMessage(messageType, p)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func main() {
	fmt.Println("Hello, World!")
	router := gin.Default()

	router.GET("/", homePage)
	router.GET("/ws", wsEndPoint)
	router.GET("/:roomName", roomMessage)

	router.Run(":8080")
}
