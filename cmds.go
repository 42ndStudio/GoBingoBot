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

	msg.Text, err = board.printText()
	if err != nil {
		strerr := fmt.Sprintf("failed printing text board (%s) for game (%s)", board.BoardID, gameID)
		logError(strerr, err)
		return msg, errors.New(strerr)
	}

	return msg, nil
}
