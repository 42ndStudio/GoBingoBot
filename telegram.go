// 42nd Studio @2020
// MuchLove
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const errorOccurred = "Ocurrió un error x.o\n"

var gifs = map[string][]string{
	"hello":    {"https://www.reactiongifs.com/r/fgwv.gif"},
	"bingo":    {"https://media.giphy.com/media/BLjbqh9Yg2LqzBIbf6/giphy.gif", "https://media.giphy.com/media/6CYXe7Hf8FZyU/giphy.gif", "https://media.giphy.com/media/J8NaR2tsCdNew/giphy.gif", "https://media.giphy.com/media/vmon3eAOp1WfK/giphy.gif", "https://media.giphy.com/media/kyLYXonQYYfwYDIeZl/giphy.gif", "https://media.giphy.com/media/3oz8xRF0v9WMAUVLNK/giphy.gif", "https://media.giphy.com/media/ZUomWFktUWpFu/giphy.gif", "https://media.giphy.com/media/3oKIPnF9GF2wT4NU1a/giphy.gif", "https://media.giphy.com/media/3o6Zt8MgUuvSbkZYWc/giphy.gif", "https://media.giphy.com/media/26vUwP2u6YZ5U6WSk/giphy.gif", "https://media.giphy.com/media/Is1O1TWV0LEJi/giphy.gif", "https://media.giphy.com/media/s2qXK8wAvkHTO/giphy.gif"},
	"fail":     {"https://media.giphy.com/media/N35rW3vRNeaDC/giphy.gif", "https://media.giphy.com/media/1VT3UNeWdijUSMpRL4/giphy.gif", "https://media.giphy.com/media/dlMIwDQAxXn1K/giphy.gif", "https://media.giphy.com/media/biKilc2r3kGOyXNiDq/giphy.gif", "https://media.giphy.com/media/Q61LJj43H48z1FIK4X/giphy-downsized-large.gif", "https://media.giphy.com/media/l1J9EdzfOSgfyueLm/giphy.gif"},
	"confused": {"https://media.giphy.com/media/gKsJUddjnpPG0/giphy.gif", "https://media.giphy.com/media/xT0BKmtQGLbumr5RCM/giphy.gif", "https://media.giphy.com/media/ji6zzUZwNIuLS/giphy.gif", "https://media.giphy.com/media/WRQBXSCnEFJIuxktnw/giphy.gif", "https://media.giphy.com/media/a0FuPjiLZev4c/giphy.gif", "https://media.giphy.com/media/xDQ3Oql1BN54c/giphy.gif", "https://media.giphy.com/media/12xsYM8AbsyoCs/giphy.gif", "https://media.giphy.com/media/nWn6ko2ygIeEU/giphy.gif", "https://media.giphy.com/media/XQvhpuryrPGnK/giphy.gif", "https://media.giphy.com/media/lkdH8FmImcGoylv3t3/giphy.gif", "https://media.giphy.com/media/jkojXEIwuqp6o/giphy.gif", "https://media.giphy.com/media/fpXxIjftmkk9y/giphy.gif", "https://media.giphy.com/media/vh9isNb4S2Spa/giphy.gif", "https://media.giphy.com/media/iHe7mA9M9SsyQ/giphy.gif"},
}

var masterKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Listar Juegos"),
		tgbotapi.NewKeyboardButton("Nuevo Juego"),
	),
)

func gimmeOrganizerKeyboard(telegramID string) *tgbotapi.ReplyKeyboardMarkup {
	var (
		mk   tgbotapi.ReplyKeyboardMarkup
		rows [][]tgbotapi.KeyboardButton
	)

	_, _, err := getGameNOrganizer(telegramID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game from organizer (TG %s)", telegramID)
		logError(strerr, err)
		return nil
	}

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Nueva partida"),
		tgbotapi.NewKeyboardButton("Cambiar Modo de Juego"),
	))

	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Vender Tablero"),
	))

	if waitingon[telegramID] == "check_game" {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Terminar Comprobacion"),
		))
	} else {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Comprobar Tablero"),
		))
	}

	mk = tgbotapi.NewReplyKeyboard(rows...)
	return &mk
}

// Arreglo donde el key es el Telegram ID de quien estamos esperando una respuesta, el valor es la operacion pendiente
var waitingon map[string]string

func randomGif(ident string) string {
	return gifs[ident][rand.Intn(len(gifs[ident]))]
}

func runTelegramUpdater() {
	print("telegram updater starting o.O ")
	time.Sleep(time.Millisecond * 2000)

	// Hard-code generate boards on start
	// for i := 0; i < 500; i++ {
	// 	cmdBoardNew(se.masterID, "VFNFTOPSNWUB3468NIL04", tgbotapi.MessageConfig{}, strconv.Itoa(i+1)+".png")
	// }

	var err error
	se.bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_KEY"))
	if err != nil {
		log.Panic(err)
	}

	se.bot.Debug = true

	log.Printf("Authorized on account %s, master is: %s", se.bot.Self.UserName, se.masterID)

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
		fromID = strconv.Itoa(update.CallbackQuery.From.ID)
	} else {
		mesOb = update.Message
		fromID = strconv.Itoa(mesOb.From.ID)
	}
	if mesOb != nil && mesOb.From != nil {
		name = mesOb.From.FirstName
	}
	msg = tgbotapi.NewMessage(mesOb.Chat.ID, "")
	comando = strings.ToLower(comando)
	switch comando {
	case "start", "hi", "hello", "hola":
		msg.Text = fmt.Sprintf("👋 Hola %s!!\n\n¿Cúal es tu código?\nWhat's your code?", name)
		waitingon[fromID] = "starting"
	case "master", "👽":
		if fromID != se.masterID {
			msg.Text = randomGif("fail")
		} else {
			msg.Text = randomGif("hello")
			msg.ReplyMarkup = masterKeyboard
		}
	case "listar juegos", "list":
		if fromID != se.masterID {
			msg.Text = randomGif("fail")
		} else {
			mk, err := gamesList()
			if err != nil {
				msg.Text = errorOccurred + randomGif("fail")
			} else {
				msg.Text = "Juegos Activos:"
				msg.ReplyMarkup = mk
			}
		}
	case "nuevo juego", "new":
		if fromID != se.masterID {
			msg.Text = randomGif("fail")
		} else {
			msg.Text = "Nombre pare le juego:"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			waitingon[fromID] = "new_game"
		}
	case "status", "/status", "❤️":
		msg.Text = "I'm ok. 💜 Love 4 U"
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "nueva partida":
		msg = cmdStartGame(mesOb, msg, name, fromID)
	case "vender tablero":
		msg = cmdSellBoard(msg, fromID)
	case "comprobar tablero":
		msg = cmdCheckBoardPre(msg, fromID)
		msg.ReplyMarkup = gimmeOrganizerKeyboard(fromID)
	case "terminar comprobacion":
		delete(waitingon, fromID)
		msg.Text = "ok, volvemos al juego, quedo atento a balotas"
		msg.ReplyMarkup = gimmeOrganizerKeyboard(fromID)
	case "cambiar modo de juego":
		msg = cmdChangeModeTrigger(mesOb, msg, fromID)
	case "generar":
		if fromID == se.masterID {
			msg = cmdGenerateBoards(msg, fromID)
		} else {
			msg.Text = randomGif("fail")
		}
	default:
		if resp, ok := pendingComm(comando, mesOb, msg, name, fromID); ok {
			msg = resp
		} else if msg, ok = specialComm(comando, fromID, msg); ok {
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

func pendingComm(comando string, mesOb *tgbotapi.Message, respmsg tgbotapi.MessageConfig, name, fromID string) (tgbotapi.MessageConfig, bool) {
	var (
		ok bool
	)
	if val, won := waitingon[fromID]; won {
		// Revisar cual es el tema pendiente (de que estamos hablando)
		switch val {
		case "starting":
			ok = true // si era pendiente
			// Check Client
			err := newOrganizer(fromID, comando, name)
			if err != nil {
				strerr := "failed newOrganizer()"
				logError(strerr, err)
			} else {
				respmsg.Text = fmt.Sprintf("Correcto, ahora estas acargo")
				respmsg.ReplyMarkup = gimmeOrganizerKeyboard(fromID)
				delete(waitingon, fromID)
			}
		case "new_game":
			ok = true // si era pendiente
			if comando == "" {
				logError("missing new_game's name", nil)
			}
			game := new(BingoGame)
			game.Name = comando
			err := game.guardar()
			if err != nil {
				strerr := "failed game.guardar()"
				logError(strerr, err)
			} else {
				respmsg.Text = fmt.Sprintf("Juego %s creado", game.BingoID)
				respmsg.ReplyMarkup = masterKeyboard
				delete(waitingon, fromID)
			}
		case "check_game":
			ok = true // si era pendiente
			if comando == "" {
				logError("missing board id", nil)
				respmsg.Text = "falta el id del tablero"
			} else {
				respmsg = cmdCheckBoard(respmsg, fromID, comando)
			}
		case "change_mode":
			ok = true // si era pendiente
			if comando == "" {
				logError("missing game mode", nil)
				respmsg.Text = "a que modo?"
			} else {
				respmsg = cmdChangeMode(comando, respmsg, fromID)
			}
		}
	}
	return respmsg, ok
}

func gamesList() (tgbotapi.InlineKeyboardMarkup, error) {
	var (
		games []BingoGame
		mk    tgbotapi.InlineKeyboardMarkup
		rows  [][]tgbotapi.InlineKeyboardButton
	)

	err := se.db.Where("1 = 1").Find(&games).Error
	if err != nil {
		logError("failed loading games @gamesList", nil)
		return mk, err
	}
	for _, game := range games {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(game.Name, "game_details:"+game.BingoID)))
	}

	mk = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return mk, nil
}

func specialComm(command, fromID string, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, bool) {
	var err error
	executed := false
	msg.Text = errorOccurred
	var regs = map[string]string{
		"game_details": `(?m)^game_details:([\w-]*)`,
		"board_new":    `(?m)^board_new:([\w-]*)`,
		"bingo_balot":  `(?m)^([BINGObingo])([\d]{1,2})`,
	}

	for cmd, reg := range regs {
		var re = regexp.MustCompile(reg)
		rs := re.FindStringSubmatch(command)
		if len(rs) > 0 {
			switch cmd {
			case "game_details":
				// Game Details
				msg, err = cmdGameDetails(fromID, rs[1], msg)
				if err != nil {
					msg.Text = errorOccurred + randomGif("fail")
					return msg, true
				}
				return msg, true

			case "board_new":
				// Generate (Sell) Board
				msg, err = cmdBoardNew(fromID, rs[1], msg, "out.png")
				if err != nil {
					msg.Text = errorOccurred + randomGif("fail")
					return msg, true
				}
				return msg, true

			case "bingo_balot":
				// Drawm Balot, check organizer and game playing
				return cmdDrawnBalot(fromID, rs[1], rs[2], msg)
			}
		}
	}

	return msg, executed
}
