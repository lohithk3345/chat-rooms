package main

import (
	"chat/internal/client"
	"log"
	"os"
)

func CollectArgs() []string {
	args := os.Args
	return args
}

func main() {
	roomId := CollectArgs()[1]
	conn, err := client.Connect(roomId)
	if err != nil {
		log.Println(err)
		return
	}
	c := client.NewClient()
	c.Run(conn)
}
