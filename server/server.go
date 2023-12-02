package server

import (
	"io"
	"log"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Handle(ws *websocket.Conn) {
	log.Println("A new client joined:", ws.RemoteAddr())

	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	// buffer := make([]byte, 2048)

	ws.WriteMessage(websocket.BinaryMessage, []byte("Welcome client!"))
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Client left")
				break
			}
			log.Fatal("Error in reading")
			continue
		}
		log.Println(string(msg))
		ws.WriteMessage(websocket.BinaryMessage, []byte(msg))
	}
}
