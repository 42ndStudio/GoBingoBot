// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
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

func (game *BingoGame) loadOrganizer(telegramID string) (BingoOrganizer, error) {
	var organizer BingoOrganizer
	err := se.db.Where("bingo_id = ?", game.BingoID).First(&organizer).Error
	return organizer, err
}
