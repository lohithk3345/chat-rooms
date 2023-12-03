package main

import (
	"chat/internal/client"
	"log"
)

func main() {
	roomId := client.CollectArgs()[1]
	conn, err := client.Connect(roomId)
	if err != nil {
		log.Println(err)
		return
	}
	c := client.NewClient()
	c.Run(conn)
}
