// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
)

func (game *BingoGame) guardar() error {
	var err error
	// Check si tiene clave
	if game.Password == "" {
		game.Password, err = GenerateRandomString(8)
		if err != nil {

		}
	}
	// Crear o Actualizar
	if game.BingoID == "" {
		// Asignar UID
		game.BingoID = UIDNew(21)
		err = se.db.Create(&game).Error
	} else {
		err = se.db.Save(&game).Error
	}
	if err != nil {
		strerr := "error guardando game object"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

func (game *BingoGame) loadFromID(bingoID string) error {
	err := se.db.Where("bingo_id = ?", bingoID).First(&game).Error
	return err
}

func (game *BingoGame) loadActiveFromPass(password string) error {
	err := se.db.Where("accepting_organizers = true AND password = ?", password).First(&game).Error
	return err
}

func (game *BingoGame) getOrganizer(telegramID string) (BingoOrganizer, error) {
	var organizer BingoOrganizer
	err := se.db.Where("bingo_id = ?", game.BingoID).First(&organizer).Error
	return organizer, err
}

func (game *BingoGame) getBoard(boardID string) (BingoBoard, error) {
	var board BingoBoard
	err := se.db.Where("bingo_id = ? AND board_id = ?", game.BingoID, boardID).First(&board).Error
	return board, err
}

func (game *BingoGame) loadBoards() error {
	err := se.db.Where("bingo_id = ?", game.BingoID).Find(&game.boards).Error
	return err
}

// drawBalot registra una balota sacada
// marca los tableros que lo tienen
func (game *BingoGame) drawBalot(letter, number string) (int, error) {
	winners := 0
	if game.DrawnBalots != "" {
		game.DrawnBalots += ","
	}
	game.DrawnBalots += letter + number
	game.Playing = true
	err := game.guardar()
	if err != nil {
		strerr := fmt.Sprintf("failed game.drawBalot (%s %s)", letter, number)
		logError(strerr, err)
		return 0, errors.New(strerr)
	}

	err = game.loadBoards()
	if err != nil {
		strerr := fmt.Sprintf("failed game.loadBoards (GID %s)", game.BingoID)
		logError(strerr, err)
		return 0, errors.New(strerr)
	}

	println("marcando tableros", len(game.boards))

	for _, board := range game.boards {
		won, err := board.markSlots(letter, number, game.CurrentMode)
		if err != nil {
			strerr := fmt.Sprintf("failed marking board (ID %s)", board.BoardID)
			logError(strerr, err)
			return winners, errors.New(strerr)
		}
		if won {
			winners++
		}
	}

	return winners, err
}
