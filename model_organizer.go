// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
)

func (organizer *BingoOrganizer) guardar() error {
	var err error
	// Crear o Actualizar
	if organizer.ID == 0 {
		err = se.db.Create(&organizer).Error
	} else {
		err = se.db.Save(&organizer).Error
	}
	if err != nil {
		strerr := "error guardando organizer object"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

func newOrganizer(telegramID, password, name string) error {
	var (
		organizer BingoOrganizer
		game      BingoGame
	)
	err := game.loadActiveFromPass(password)
	if err != nil {
		strerr := "failed to loadActiveFromPass"
		logError(strerr, err)
		return errors.New(strerr)
	}

	organizer.BingoID = game.BingoID
	organizer.Name = name
	organizer.TelegramID = telegramID
	err = organizer.guardar()

	if err != nil {
		strerr := "failed organizer.guardar()"
		logError(strerr, err)
		return errors.New(strerr)
	}

	return nil
}
