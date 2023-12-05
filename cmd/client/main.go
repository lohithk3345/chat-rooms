package main

import (
	"chat/internal/client"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CollectArgs() []string {
	args := os.Args
	return args
}

var rootCommand = &cobra.Command{
	Use:     "create",
	Short:   "Create a room and join",
	Example: "go run ./cmd/client create",
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c", "cr"},
	Short:   "Create's room and joins it",
	Example: "go run ./cmd/client create",
	Run: func(cmd *cobra.Command, args []string) {

		roomId, err := client.CreateRoom()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		log.Println(roomId)

		conn, err := client.Connect(roomId)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		c := client.NewClient()
		c.Run(conn)
	},
}

var joinCmd = &cobra.Command{
	Use:     "join",
	Aliases: []string{"j"},
	Short:   "Join a room with id you specify",
	Example: "client join -r [roomId]",
	Run: func(cmd *cobra.Command, args []string) {
		roomId := viper.GetString("join")
		log.Println("RoomId", roomId)
		conn, err := client.Connect(roomId)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		c := client.NewClient()
		c.Run(conn)
	},
}

func init() {
	joinCmd.PersistentFlags().StringP("join", "r", "", "Join the room you specify")
	viper.BindPFlag("join", joinCmd.PersistentFlags().Lookup("join"))

	rootCommand.AddCommand(createCmd, joinCmd)
}

func main() {
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
