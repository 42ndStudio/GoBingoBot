// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
)

func (slot *BingoSlot) guardar() error {
	var err error
	// Crear o Actualizar
	if slot.ID == 0 {
		err = se.db.Create(&slot).Error
	} else {
		err = se.db.Save(&slot).Error
	}
	if err != nil {
		strerr := "error guardando slot object"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

func (slot *BingoSlot) loadFromBoard(bingoID, boardID, letter string, number int) error {
	err := se.db.Where("bingo_id = ? AND board_id = ? AND letter = ? AND number = ?", bingoID, boardID, letter, number).First(&slot).Error
	return err
}
