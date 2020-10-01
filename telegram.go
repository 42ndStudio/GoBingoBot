package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var gifs = map[string][]string{
	"hello":    {"https://www.reactiongifs.com/r/fgwv.gif"},
	"bingo":    {"https://media.giphy.com/media/C8CZu2szOWVuGB8dB9/giphy.gif"},
	"fail":     {"https://media.giphy.com/media/N35rW3vRNeaDC/giphy.gif", "https://media.giphy.com/media/1VT3UNeWdijUSMpRL4/giphy.gif", "https://media.giphy.com/media/dlMIwDQAxXn1K/giphy.gif", "https://media.giphy.com/media/biKilc2r3kGOyXNiDq/giphy.gif", "https://media.giphy.com/media/Q61LJj43H48z1FIK4X/giphy-downsized-large.gif", "https://media.giphy.com/media/l1J9EdzfOSgfyueLm/giphy.gif", "https://media.giphy.com/media/jTkAOfSysPWVk3r68Q/giphy.gif", "https://media.giphy.com/media/l0DELAq175S4Zeobm/giphy.gif"},
	"confused": {"https://media.giphy.com/media/gKsJUddjnpPG0/giphy.gif", "https://media.giphy.com/media/xT0BKmtQGLbumr5RCM/giphy.gif", "https://media.giphy.com/media/ji6zzUZwNIuLS/giphy.gif", "https://media.giphy.com/media/WRQBXSCnEFJIuxktnw/giphy.gif", "https://media.giphy.com/media/a0FuPjiLZev4c/giphy.gif", "https://media.giphy.com/media/xDQ3Oql1BN54c/giphy.gif", "https://media.giphy.com/media/12xsYM8AbsyoCs/giphy.gif", "https://media.giphy.com/media/nWn6ko2ygIeEU/giphy.gif", "https://media.giphy.com/media/XQvhpuryrPGnK/giphy.gif", "https://media.giphy.com/media/lkdH8FmImcGoylv3t3/giphy.gif", "https://media.giphy.com/media/jkojXEIwuqp6o/giphy.gif", "https://media.giphy.com/media/fpXxIjftmkk9y/giphy.gif", "https://media.giphy.com/media/vh9isNb4S2Spa/giphy.gif", "https://media.giphy.com/media/iHe7mA9M9SsyQ/giphy.gif"},
}

var masterKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Iniciar Juego"),
		tgbotapi.NewKeyboardButton("Comprobar Tablero"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Cambiar Modo de Juego"),
		tgbotapi.NewKeyboardButton("Ajustes Notificaciones"),
	),
)

// Arreglo donde el key es el Telegram ID de quien estamos esperando una respuesta, el valor es la operacion pendiente
var waitingon map[string]string

func randomGif(ident string) string {
	return gifs[ident][rand.Intn(len(gifs[ident]))]
}

func runTelegramUpdater() {
	print("telegram updater starting o.O ")
	var err error
	se.bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_KEY"))
	if err != nil {
		log.Panic(err)
	}

	se.bot.Debug = true

	log.Printf("Authorized on account %s", se.bot.Self.UserName)

	waitingon = make(map[string]string)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := se.bot.GetUpdatesChan(u)

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {

		if update.Message == nil && (update.CallbackQuery == nil || update.CallbackQuery.Data == "") {
			continue
		}

		if update.Message == nil && update.CallbackQuery.Data != "" {
			log.Println("Hay callback query", update.CallbackQuery.Data)
			log.Println(update.CallbackQuery.Message.Chat.ID)
			processCommand(update.CallbackQuery.Data, update)
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				log.Println("comando recibido: ", update.Message.Command())
				processCommand(update.Message.Command(), update)
			} else {
				log.Println("mensaje recibido: ", update.Message.Text)
				processCommand(update.Message.Text, update)
			}
		}
	}
}

func processCommand(comando string, update tgbotapi.Update) {
	var (
		msg          tgbotapi.MessageConfig
		mesOb        *tgbotapi.Message
		name, fromID string
	)
	if update.Message == nil {
		mesOb = update.CallbackQuery.Message
	} else {
		mesOb = update.Message
	}
	if mesOb != nil && mesOb.From != nil {
		name = mesOb.From.FirstName
		fromID = strconv.Itoa(mesOb.From.ID)
	}
	msg = tgbotapi.NewMessage(mesOb.Chat.ID, "")
	comando = strings.ToLower(comando)
	switch comando {
	case "start", "hi", "hello", "hola":
		msg.Text = fmt.Sprintf("👋 Hola %s!!\n\n¿Cúal es tu código?\nWhat's your code?", name)
		waitingon[fromID] = "starting"
	case "status", "/status", "❤️":
		msg.Text = "I'm ok. 💜 Love 4 U"
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	default:
		if resp, ok := pendingComm(mesOb, msg, name, fromID); ok {
			msg = resp
		} else {
			msg.Text = "Lo siento, no reconozco ese comando\nSorry, I don't know that command \n" + randomGif("confused")
			log.Println(fmt.Sprintf("ERR: unknown command %s", comando))
		}
	}
	se.bot.Send(msg)
}

func telegramMessage(toID int64, message string) tgbotapi.MessageConfig {
	var (
		msg tgbotapi.MessageConfig
	)

	msg = tgbotapi.NewMessage(toID, message)

	return msg
}

func pendingComm(mesOb *tgbotapi.Message, respmsg tgbotapi.MessageConfig, name, fromID string) (tgbotapi.MessageConfig, bool) {
	var (
		ok bool
	)
	if val, won := waitingon[fromID]; won {
		// Revisar cual es el tema pendiente (de que estamos hablando)
		switch val {
		case "starting":
			// Check Client
			delete(waitingon, fromID)
			ok = true // si era pendiente
		}
	}
	return respmsg, ok
}
