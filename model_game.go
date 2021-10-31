// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

var checkAllBoards bool

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
	fmt.Println("loading game", bingoID)
	err := se.db.Where("bingo_id = ?", strings.ToUpper(bingoID)).First(&game).Error
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

func (game *BingoGame) getBoardByInt(boardID int64) (BingoBoard, error) {
	var board BingoBoard
	err := se.db.Where("bingo_id = ? AND id = ?", game.BingoID, boardID).First(&board).Error
	return board, err
}

func (game *BingoGame) loadBoards() error {
	err := se.db.Where("bingo_id = ?", game.BingoID).Find(&game.boards).Error
	return err
}

func (game *BingoGame) isUnique(hash string) (bool, error) {
	var existingBoards []BingoBoard
	println("checky", game.BingoID, hash)
	err := se.db.Where("bingo_id = ? AND board_hash = ?", game.BingoID, hash).Find(&existingBoards).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		strerr := fmt.Sprintf("failed checking board uniqueness gid: %s", game.BingoID)
		logError(strerr, err)
		return false, errors.New(strerr)
	}
	if len(existingBoards) > 0 {
		logError("non unique board!!!", err)
		return false, nil
	}
	return true, nil
}

// drawBalot registra una balota sacada
// marca los tableros que lo tienen
func (game *BingoGame) drawBalot(letter, number string) (int, error) {
	winners := 0

	if stringInSlice(letter+number, strings.Split(game.DrawnBalots, ",")) {
		strerr := fmt.Sprintf("already drawn (%s %s)", letter, number)
		logError(strerr, nil)
		return -42, errors.New(strerr)
	}

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

	if !checkAllBoards {
		return 0, nil
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
