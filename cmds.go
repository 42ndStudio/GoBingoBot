// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var boardIDMode = "int"

const ONLINE = false

func cmdGameDetails(fromID, gameID string, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	fmt.Println(gameID)
	log.Println("Detalles de juego", gameID, "solicitados")
	game, _, err := getGameNOrganizer(fromID, gameID)
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer (TG %s) from id for details", fromID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}

	var (
		rows [][]tgbotapi.InlineKeyboardButton
	)
	msg.ParseMode = "markdown"
	msg.Text = "**Juego:** " + game.Name + "\n" +
		"**ID:** " + game.BingoID + "\n" +
		"**Vendidos:** " + strconv.Itoa(game.BoardsSold) + "\n" +
		"**Abierto:** " + strconv.FormatBool(game.AcceptingOrganizers) + "\n"
	if game.AcceptingOrganizers {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("No mas Organizadores", "organizers_no:"+game.BingoID))))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("Aceptar Organizadores", "organizers_ok:"+game.BingoID))))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("Vender Tablero", "board_new:"+game.BingoID))))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("Cambiar Dinamica", "game_dynamic_new:"+game.BingoID))))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("Listar Oranizadores", "game_organizers:"+game.BingoID))))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return msg, nil
}

func cmdBoardNew(fromID, gameID string, msg tgbotapi.MessageConfig, outName string) (tgbotapi.MessageConfig, error) {
	fmt.Println(gameID)
	log.Println("Nuevo tablero de juego", gameID, "solicitado por", fromID)
	game, organizer, err := getGameNOrganizer(fromID, gameID)
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer (TG %s) from id for details", gameID, fromID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}

	var board BingoBoard
	board.BingoID = game.BingoID
	board.Sold = true
	err = board.guardar()
	if err != nil {
		strerr := fmt.Sprintf("failed saving board for game (%s)", gameID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}
	err = board.generate(game)
	if err != nil {
		strerr := fmt.Sprintf("failed generating board for game (%s)", gameID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}

	if organizer != nil {
		organizer.BoardsSold++
		err = organizer.guardar()
		if err != nil {
			strerr := fmt.Sprintf("failed incrementing organizer (%d) sold boards", organizer.ID)
			logError(strerr, err)
		}
	}

	game.BoardsSold++
	err = game.guardar()
	if err != nil {
		strerr := fmt.Sprintf("failed incrementing game (%d) sold boards", organizer.ID)
		logError(strerr, err)
	}

	go func() {
		err = board.drawImage(outName)
		if err != nil {
			strerr := fmt.Sprintf("failed drawing board (%s) for game (%s)", board.BoardID, gameID)
			logError(strerr, err)
			return
		}

		if ONLINE {
			imsg := tgbotapi.NewPhotoUpload(msg.ChatID, outName)
			imsg.Caption = board.BoardID

			se.bot.Send(imsg)

			go func() {
				msg, err = cmdGameDetails(fromID, game.BingoID, msg)
				if err == nil {
					se.bot.Send(msg)
				}
			}()
		}
	}()

	msg.Text = "cargando..."
	//if err != nil {
	//	strerr := fmt.Sprintf("failed printing text board (%s) for game (%s)", board.BoardID, gameID)
	//	logError(strerr, err)
	//	return msg, errors.New(strerr)
	//}

	return msg, nil
}

// cmdSellBoard ejecutado por organizador para vender (generar) un tablero
func cmdSellBoard(respmsg tgbotapi.MessageConfig, fromID string) tgbotapi.MessageConfig {
	var ()
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game from organizer (TG %s) @sellBoard", fromID)
		logError(strerr, err)
		return respmsg
	}

	respmsg, err = cmdBoardNew(fromID, game.BingoID, respmsg, "out.png")
	if err != nil {
		respmsg.Text = errorOccurred + randomGif("fail")
	}

	return respmsg
}

// cmdStartGame ejecutado por organizador para iniciar una nueva partida
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdStartGame(mesOb *tgbotapi.Message, respmsg tgbotapi.MessageConfig, name, fromID string) tgbotapi.MessageConfig {

	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	ogdrawn := game.DrawnBalots
	game.DrawnBalots = ""
	game.Playing = true
	err = game.guardar()
	if err != nil {
		strerr := fmt.Sprintf("failed saving game.started %s", game.BingoID)
		logError(strerr, err)
		return respmsg
	}

	err = game.loadBoards()
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) boards", game.BingoID)
		logError(strerr, err)
		return respmsg
	}

	for _, board := range game.boards {
		board.clearSlots()
	}

	delete(waitingon, fromID)

	textP1 := ""
	for _, drawn := range strings.Split(ogdrawn, ",") {
		textP1 += drawn + "\n"
	}

	respmsg.Text = "Balotas sacada: \n" + textP1 + "\n\nNueva  Partida Iniciada!\n GoBinGO! 😜\n\n Enviame un mensaje con cada balota que saques"

	return respmsg
}

// cmdDrawnBalot ejecutado por organizador para iniciar una nueva partida
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdDrawnBalot(fromID, letter, number string, respmsg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, bool) {
	var ()
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg, false
	}

	if !game.Playing {
		strerr := fmt.Sprintf("drawn balot but game is not playing (%s)", game.BingoID)
		respmsg.Text = "No se esta jugando una partida! Inicia una nueva partida para jugar"
		logError(strerr, err)
		return respmsg, false
	}

	winners, err := game.drawBalot(letter, number)
	if err != nil {
		strerr := fmt.Sprintf("failed game.drawBalot %s", game.BingoID)
		logError(strerr, err)
		if winners == -42 {
			respmsg.Text = randomGif("fail") + "\n Ya habia sido registrado"
		}
		return respmsg, true
	}

	respmsg.Text = fmt.Sprintf("Registrado (%s-%s)", letter, number)
	if number == "42" {
		respmsg.Text += "(😛 woooo 42nd!!)"
	}
	if winners > 0 {
		respmsg.Text += fmt.Sprintf("\n\nHabemus ganadores: %d", winners)
	}

	return respmsg, true
}

// cmdCheckBoardPre pone a la espera chequeo de tablero
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdCheckBoardPre(respmsg tgbotapi.MessageConfig, fromID string) tgbotapi.MessageConfig {
	var ()
	respmsg.Text = errorOccurred + randomGif("fail")

	_, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	waitingon[fromID] = "check_game"

	respmsg.Text = "Que tablero quieres comprobar?"

	return respmsg
}

// cmdCheckBoard chequea tablero
func cmdCheckBoard(respmsg tgbotapi.MessageConfig, fromID, boardID string) tgbotapi.MessageConfig {
	var (
		board BingoBoard
	)
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	if boardIDMode == "int" {
		intid, _ := strconv.ParseInt(boardID, 10, 64)
		if err != nil {
			strerr := fmt.Sprintf("failed parsing board integer id %s", boardID)
			logError(strerr, err)
			return respmsg
		}
		board, err = game.getBoardByInt(intid)
		if err != nil {
			strerr := fmt.Sprintf("failed getting board %s", boardID)
			logError(strerr, err)
			return respmsg
		}
	} else {
		board, err = game.getBoard(boardID)
		if err != nil {
			strerr := fmt.Sprintf("failed getting board %s", boardID)
			logError(strerr, err)
			return respmsg
		}
	}

	drawn := strings.Split(game.DrawnBalots, ",")

	var winner bool
	for _, draw := range drawn {
		winner, err = board.markSlots(string(draw[0]), string(draw[1:]), game.CurrentMode)
		if err != nil {
			strerr := fmt.Sprintf("failed checking board (%s) at game (%s)", board.BoardID, game.BingoID)
			logError(strerr, err)
			return respmsg
		}
	}

	if winner {
		respmsg.Text = fmt.Sprintf("🥳 BINGO!\nEl tablero %s ES ganador!!!", boardID) + "\n" + randomGif("bingo")
		println("is winner")
	} else {
		respmsg.Text = fmt.Sprintf("El tablero %s NO es ganador 😞", boardID)
		println("is loser")
	}

	return respmsg
}

// cmdChangeModeTrigger cambia el modo de juego (si es organizador)
func cmdChangeModeTrigger(inmsg *tgbotapi.Message, respmsg tgbotapi.MessageConfig, fromID string) tgbotapi.MessageConfig {
	respmsg.Text = errorOccurred + randomGif("fail")

	_, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	waitingon[fromID] = "change_mode"

	respmsg.Text = "Que modo desean jugar?"

	var (
		rows [][]tgbotapi.InlineKeyboardButton
	)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow((tgbotapi.NewInlineKeyboardButtonData("Completo", "a"))))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Figura 'O'", "o"),
		tgbotapi.NewInlineKeyboardButtonData("'N'", "n"),
		tgbotapi.NewInlineKeyboardButtonData("'C'", "c"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("'/'", "/"),
		tgbotapi.NewInlineKeyboardButtonData("'\\'", "\\"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Col 'B'", "lb"),
		tgbotapi.NewInlineKeyboardButtonData("'I'", "li"),
		tgbotapi.NewInlineKeyboardButtonData("'N'", "ln"),
		tgbotapi.NewInlineKeyboardButtonData("'G'", "lg"),
		tgbotapi.NewInlineKeyboardButtonData("'O'", "lo"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Fil '1'", "l1"),
		tgbotapi.NewInlineKeyboardButtonData("'2'", "l2"),
		tgbotapi.NewInlineKeyboardButtonData("'3'", "l3"),
		tgbotapi.NewInlineKeyboardButtonData("'4'", "l4"),
		tgbotapi.NewInlineKeyboardButtonData("'5'", "l5"),
	))
	respmsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return respmsg
}

// cmdChangeMode cambia el modo de juego (si es organizador)
func cmdChangeMode(gameMode string, respmsg tgbotapi.MessageConfig, fromID string) tgbotapi.MessageConfig {
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	modes := []string{"lb", "li", "ln", "lg", "lo", "a", "c", "o", "n", "/", "\\", "l1", "l2", "l3", "l4", "l5"}

	println("changing to mode " + gameMode)
	println(fmt.Sprint(modes))

	if stringInSlice(gameMode, modes) {
		respmsg.Text = "Modo de juego cambio a: " + gameMode
		game.CurrentMode = gameMode
		err = game.guardar()
		if err != nil {
			strerr := "failed saving game new mode"
			logError(strerr, err)
			return respmsg
		}
		delete(waitingon, fromID)
	} else {
		respmsg.Text = "Modo invalido: " + gameMode
	}

	return respmsg
}

// cmdGenerateBoards
func cmdGenerateBoards(respmsg tgbotapi.MessageConfig, fromID string) tgbotapi.MessageConfig {
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game and organizer from (TG %s) for details", fromID)
		logError(strerr, err)
		return respmsg
	}

	for i := 0; i < 500; i++ {
		respmsg, err = cmdBoardNew(fromID, game.BingoID, respmsg, strconv.Itoa(i+1)+".png")
		if err != nil {
			respmsg.Text = errorOccurred + randomGif("fail")
		}

	}

	respmsg.Text = "GoBinGO! 😜\n\n Tableros Generados!"

	return respmsg
}
