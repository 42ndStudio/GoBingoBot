package main

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/jinzhu/gorm"
)

var (
	se server
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	// ID del usuario maestro
	masterID int64

	bot *tgbotapi.BotAPI

	db *gorm.DB
}

func main() {
	println("\033]0;GoBingo\007")

	var err error

	se.masterID, err = strconv.ParseInt(os.Getenv("MASTER_ID"), 10, 64)
	if err != nil {
		panic("Missing env['master_id']")
	}

	go runTelegramUpdater()

	select {}
}
