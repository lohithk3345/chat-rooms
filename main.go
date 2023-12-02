package main

import (
	"chat/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	// s := server.NewServer()
	router := gin.Default()
	router.GET("/createRoom", handleCreateRoom)
	router.GET("/ws/:id", func(ctx *gin.Context) { handleWS(ctx) })
	// router.GET("/join:room")
	http.ListenAndServe(":3000", router)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}
)

var rooms_map = make(map[string]*server.Room)
var clients = make(map[string]*server.Client)

func handleWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error in upgrading")
		return
	}

	room_id := ctx.Param("id")
	log.Println(room_id)
	room := rooms_map[room_id]

	conn.WriteMessage(websocket.TextMessage, []byte(room_id))

	client, id := server.NewClient(room, conn)
	// s.Handle(conn)
	clients[id] = client

	client.JoinRoom(room)

	go client.Read()
	go client.Write()

	select {}
}

func handleCreateRoom(c *gin.Context) {
	room, id := server.CreateRoom()
	rooms_map[id] = room
	res := fmt.Sprintf("Room Created id: %s\n", id)
	c.Writer.Write([]byte(res))
	return
}
