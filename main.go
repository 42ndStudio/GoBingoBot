// With love
// 42nd Studio
// 2020-2021

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

	go startSQLLiteDB()
	go runWebServer()

	se.masterID = os.Getenv("MASTER_ID")
	if se.masterID != "" {
		if err != nil {
			panic("Missing env['master_id']")
		}
		go runTelegramUpdater()
	}

	select {}
}
