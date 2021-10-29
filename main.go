// 42nd Studio @2020
// MuchLove
package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/jinzhu/gorm"
)

var (
	se server
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	// ID del usuario maestro
	masterID string

	bot *tgbotapi.BotAPI

	db *gorm.DB
}

func main() {
	println("\033]0;GoBingo\007")

	var err error

	se.masterID = os.Getenv("MASTER_ID")
	if err != nil {
		panic("Missing env['master_id']")
	}

	go startSQLLiteDB()
	go runTelegramUpdater()

	select {}
}
