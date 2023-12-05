package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var EnvMap = getENV()

func isEmpty(str string) bool {
	return str == ""
}

func getENV() map[string]string {
	envMap := make(map[string]string)
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
	envMap["HOST_URL"] = os.Getenv("HOST_URL")
	envMap["HOST_PORT"] = os.Getenv("HOST_PORT")
	envMap["WEBSOCKET_URL"] = os.Getenv("WEBSOCKET_URL")
	log.Println(envMap)
	for key, value := range envMap {
		if isEmpty(value) {
			log.Fatalf("ENV not found for %s", key)
		}
	}
	return envMap
}
