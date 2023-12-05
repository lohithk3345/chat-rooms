package client

import (
	"bufio"
	"chat/utils/constants"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu     sync.Mutex
	Done   chan struct{}
	Sender chan []byte
}

func CreateRoom() (string, error) {
	uri := fmt.Sprintf("%s/createRoom", constants.HOST_URL)
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "nil", err
	}
	type CreateResponse struct {
		Id string `json:"id"`
	}
	var createResponse CreateResponse
	json.Unmarshal(resp, &createResponse)
	return createResponse.Id, nil
}

func NewClient() *Client {
	return &Client{
		Done:   make(chan struct{}),
		Sender: make(chan []byte, 2048),
	}
}

func (c *Client) SendMessage(conn *websocket.Conn, msg []byte) {
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func (c *Client) ReadMessage(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Received> ", string(msg))
		fmt.Print("Input>")
	}
}

func Connect(roomId string) (*websocket.Conn, error) {
	url := fmt.Sprintf("%s/%s", constants.WEBSOCKET_URL, roomId)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Client) userInput() {
	defer close(c.Done)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Input> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "exit" || input == "q" {
			c.Done <- struct{}{}
			return
		}
		c.mu.Lock()
		c.Sender <- []byte(input)
		c.mu.Unlock()
	}
}

func (c *Client) Run(conn *websocket.Conn) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	defer close(c.Sender)
	defer conn.Close()

	go c.ReadMessage(conn)
	go c.userInput()

	for {
		select {
		case <-c.Done:
			return
		case msg := <-c.Sender:
			c.SendMessage(conn, msg)
		case <-interrupt:
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}

}
