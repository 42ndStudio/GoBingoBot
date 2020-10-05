// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	return se.db.Where("board_id = ?", boardID).First(&board).Error
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
func (board *BingoBoard) generateSlot(x, y int) error {
	var (
		slot        BingoSlot
		printedNums []int
		letters     = []string{"B", "I", "N", "G", "O"}
	)

	println("generating slot letter:", letters[x])
	slot.BoardID = board.BoardID
	slot.Letter = letters[x]

	println("loaded slots: ", len(board.slots))

	for _, eslot := range board.slots {
		if eslot.Letter == slot.Letter {
			printedNums = append(printedNums, eslot.Number)
		}
	}

	for i := 0; i < 100; i++ {
		slot.Number = rand.Intn(((x+1)*15)-((x*15)+1)) + ((x * 15) + 1)
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

	slot.Y = y
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
		for r := 0; r < 5; r++ {
			board.generateSlot(l, r)
		}
	}

	return nil
}

func (board *BingoBoard) clearSlots() error {
	return se.db.Table("bingo_slots").Where("board_id = ?", board.BoardID).Update("marked", false).Error
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
		msg += printSlot(rows["B"][r]) + " " + strconv.Itoa(rows["I"][r]) + " " + strconv.Itoa(rows["N"][r]) + " " + strconv.Itoa(rows["G"][r]) + " " + strconv.Itoa(rows["O"][r]) + "\n"
	}

	return msg, nil
}

// markSlot marca una casilla (si la tiene)
// retorana si es ganador de la actual dinamica
func (board *BingoBoard) markSlots(letter, number, dinamica string) (bool, error) {
	var (
		haveInPlace int
		winner      bool
	)

	dinamica = strings.ToLower(dinamica)
	err := board.loadSlots()
	if len(board.slots) == 0 {
		strerr := "tablero no tiene slots!"
		logError(strerr, err)
		return false, errors.New(strerr)
	}

	numberInt, _ := strconv.Atoi(number)
	println("marcando slots", letter, numberInt)

	needTogether := needFromDynamic(dinamica)

	for _, slot := range board.slots {
		lowerLetter := strings.ToLower(slot.Letter)
		if lowerLetter == strings.ToLower(letter) && slot.Number == numberInt {
			slot.Marked = true
			err = slot.guardar()
			if err != nil {
				logError(fmt.Sprintf("Fallo marcando slot %s %s en tablero %s", letter, number, board.BoardID), err)
			}
		}

		if slot.Marked {
			if dinamica == "\\" {
				x := letter2X(lowerLetter)
				if x == slot.Y {
					haveInPlace++
				}
			} else if dinamica == "/" {
				x := letter2X(lowerLetter)
				if (x-4)*-1 == slot.Y {
					haveInPlace++
				}
			} else if dinamica == "a" {
				haveInPlace++
			} else if dinamica == "o" || dinamica == "c" {
				if slot.Y == 0 || slot.Y == 4 {
					haveInPlace++
				} else if lowerLetter == "b" || (dinamica == "o" && lowerLetter == "o") {
					haveInPlace++
				}
			} else if dinamica[0] == 'l' {
				if string(dinamica[1]) == lowerLetter {
					haveInPlace++
				} else {
					linea, err := strconv.Atoi(string(dinamica[1]))
					if err == nil && linea > 0 && linea <= 5 && slot.Y == linea-1 {
						haveInPlace++
					}
				}
			}
		}
	}

	println(fmt.Sprintf("need %d have %d", needTogether, haveInPlace))
	if haveInPlace >= needTogether {
		winner = true
	}

	return winner, nil
}
