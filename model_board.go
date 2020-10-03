// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

func (board *BingoBoard) guardar() error {
	var err error
	// Crear o Actualizar
	if board.BoardID == "" {
		// Asignar UID
		board.BoardID = UIDNew(12)
		err = se.db.Create(&board).Error
	} else {
		err = se.db.Save(&board).Error
	}
	if err != nil {
		strerr := "error guardando board object"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

func (board *BingoBoard) loadFromID(boardID string) error {
	err := se.db.Where("board_id = ?", boardID).First(&board).Error
	return err
}

// loadSlots carga las casillas del juego
func (board *BingoBoard) loadSlots() error {
	err := se.db.Where("board_id = ?", board.BoardID).Find(&board.slots).Error
	if err != nil {
		strerr := "error cargando slots de tablero"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

// generateSlot genera un slot del tablero de juego
func (board *BingoBoard) generateSlot(l int) error {
	var (
		slot        BingoSlot
		printedNums []int
		letters     = []string{"B", "I", "N", "G", "O"}
	)

	println("generating slot letter:", letters[l])
	slot.BoardID = board.BoardID
	slot.Letter = letters[l]

	println("loaded slots: ", len(board.slots))

	for _, eslot := range board.slots {
		if eslot.Letter == slot.Letter {
			printedNums = append(printedNums, eslot.Number)
		}
	}

	for i := 0; i < 100; i++ {
		slot.Number = rand.Intn(((l+1)*15)-((l*15)+1)) + ((l * 15) + 1)
		unique := true
		println(fmt.Sprintf("printedNums %v", printedNums))
		for _, pn := range printedNums {
			if pn == slot.Number {
				unique = false
				println("not unique!", pn)
				break
			}
		}
		if unique {
			break
		}
		println("will try another num")
		if i == 99 {
			strerr := "failed to generate unique number in 100 attempts"
			logError(strerr, nil)
			return errors.New(strerr)
		}
	}

	err := slot.guardar()
	if err != nil {
		strerr := "error guardando slot"
		logError(strerr, err)
		return errors.New(strerr)
	}

	board.slots = append(board.slots, slot)

	return nil
}

// generate genera el tablero de juego
func (board *BingoBoard) generate() error {
	err := board.loadSlots()
	if len(board.slots) > 0 {
		strerr := "tablero ya tiene slots generados"
		logError(strerr, err)
		return errors.New(strerr)
	}

	for l := 0; l < 5; l++ {
		for c := 0; c < 5; c++ {
			board.generateSlot(l)
		}
	}

	return nil
}

// printText devuelve el tablero impreso en texto
func (board *BingoBoard) printText() (string, error) {
	err := board.loadSlots()
	if err != nil {
		strerr := "error cargando slots @board.printText"
		logError(strerr, err)
		return "", errors.New(strerr)
	}

	msg := fmt.Sprintf("**Tablero %s**\n  B  I  N  G  O  \n", board.BoardID)

	rows := make(map[string][]int, 5)

	for _, slot := range board.slots {
		rows[slot.Letter] = append(rows[slot.Letter], slot.Number)
	}

	for r := 0; r < 5; r++ {
		msg += strconv.Itoa(rows["B"][r]) + " " + strconv.Itoa(rows["I"][r]) + " " + strconv.Itoa(rows["N"][r]) + " " + strconv.Itoa(rows["G"][r]) + " " + strconv.Itoa(rows["B"][r]) + " " + strconv.Itoa(rows["O"][r]) + "\n"
	}

	return msg, nil
}
