// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func cmdGameDetails(fromID, gameID string, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	fmt.Println(gameID)
	log.Println("Detalles de juego", gameID, "solicitados")
	game, _, err := getGameNOrganizer(fromID, gameID)
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer (TG %s) from id for details", gameID, fromID)
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

func cmdBoardNew(fromID, gameID string, msg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, error) {
	fmt.Println(gameID)
	log.Println("Nuevo tablero de juego", gameID, "solicitado")
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
	err = board.generate()
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

	msg.Text, err = board.printText()
	if err != nil {
		strerr := fmt.Sprintf("failed printing text board (%s) for game (%s)", board.BoardID, gameID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}

	return msg, nil
}

// cmdStartGame ejecutado por organizador para iniciar una nueva partida
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdStartGame(mesOb *tgbotapi.Message, respmsg tgbotapi.MessageConfig, name, fromID string) tgbotapi.MessageConfig {
	var ()
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer from (TG %s) for details", game.BingoID, fromID)
		logError(strerr, err)
		return respmsg
	}

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

	respmsg.Text = "Partida Iniciada!\n GoBinGO! 😜\n\n Enviame un mensaje con cada balota que saques"

	return respmsg
}

// cmdDrawnBalot ejecutado por organizador para iniciar una nueva partida
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdDrawnBalot(fromID, letter, number string, respmsg tgbotapi.MessageConfig) (tgbotapi.MessageConfig, bool) {
	var ()
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer from (TG %s) for details", game.BingoID, fromID)
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

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer from (TG %s) for details", game.BingoID, fromID)
		logError(strerr, err)
		return respmsg
	}

	waitingon[fromID] = "check_game"

	respmsg.Text = "Que tablero quieres comprobar?"

	return respmsg
}

// cmdCheckBoard pone a la espera chequeo de tablero
// Al iniciar una nueva partida los mensajes que  recibamos los intentaremos interpretar como balotas (ej: B5)
func cmdCheckBoard(respmsg tgbotapi.MessageConfig, fromID, boardID string) tgbotapi.MessageConfig {
	var (
		board BingoBoard
	)
	respmsg.Text = errorOccurred + randomGif("fail")

	game, _, err := getGameNOrganizer(fromID, "")
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) and organizer from (TG %s) for details", game.BingoID, fromID)
		logError(strerr, err)
		return respmsg
	}

	board, err = game.getBoard(boardID)
	if err != nil {
		strerr := fmt.Sprintf("failed getting board %s", boardID)
		logError(strerr, err)
		return respmsg
	}

	winner, err := board.markSlots("", "", game.CurrentMode)
	if err != nil {
		strerr := fmt.Sprintf("failed checking board (%s) at game (%s)", board.BoardID, game.BingoID)
		logError(strerr, err)
		return respmsg
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
