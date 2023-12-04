package main

import (
	"chat/internal/server"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	router := gin.Default()
	router.GET("/createRoom", handleCreateRoom)
	router.GET("/ws/:id", func(ctx *gin.Context) { handleWS(ctx) })
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

type ErrorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func handleWS(ctx *gin.Context) {
	room_id := ctx.Param("id")
	log.Println(room_id)
	room, isPresent := rooms_map[room_id]

	if !isPresent {
		errorJSON, err := json.Marshal(ErrorMessage{Status: 400, Message: "No Room Found"})
		if err != nil {
			log.Println(err)
			return
		}
		ctx.AbortWithStatusJSON(400, errorJSON)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Error in upgrading")
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte(room_id))

	client, id := server.NewClient(room, conn)
	clients[id] = client

	client.JoinRoom(room)

	go client.Read()
	go client.Write()
}

func handleCreateRoom(c *gin.Context) {
	room, id := server.CreateRoom()
	rooms_map[id] = room
	// res := fmt.Sprintf("Room Created id: %s\n", id)
	type CreateResponse struct {
		Id string `json:"id"`
	}
	res, _ := json.Marshal(&CreateResponse{Id: id})
	c.Writer.Write([]byte(res))
	return
}
