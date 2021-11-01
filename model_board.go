// 42nd Studio @2020
// MuchLove
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

func (board *BingoBoard) guardar() error {
	var err error
	// Crear o Actualizar
	if board.BoardID == "" {
		// Asignar UID
		board.BoardID = UIDNew(6)
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

func (board *BingoBoard) loadFromID(bingoID, boardID string) error {
	return se.db.Where("bingo_id = ? AND board_id = ?", bingoID, boardID).First(&board).Error
}

// loadSlots carga las casillas del juego
func (board *BingoBoard) loadSlots() error {
	err := se.db.Where("bingo_id = ? AND board_id = ?", board.BingoID, board.BoardID).Find(&board.slots).Error
	if err != nil {
		strerr := "error cargando slots de tablero"
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

// generateSlot genera un slot del tablero de juego
func (board *BingoBoard) generateSlot(x, y int, marked bool) (string, error) {
	var (
		slot        BingoSlot
		printedNums []int
		letters     = []string{"B", "I", "N", "G", "O"}
	)

	println("generating slot letter:", letters[x])
	slot.BingoID = board.BingoID
	slot.BoardID = board.BoardID
	slot.Letter = letters[x]
	slot.Marked = marked

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
			return "", errors.New(strerr)
		}
	}

	slot.Y = y
	err := slot.guardar()
	if err != nil {
		strerr := "error guardando slot"
		logError(strerr, err)
		return "", errors.New(strerr)
	}

	board.slots = append(board.slots, slot)

	return slot.Letter + strconv.Itoa(slot.Number), nil
}

// generate genera el tablero de juego
func (board *BingoBoard) generate(game *BingoGame) error {
	var md string
	unique := true
	for i := 0; i < 100; i++ {
		println(fmt.Sprintf("\033[33m genBoard iteration %d \033[39m", i))
		err := board.loadSlots()
		if len(board.slots) > 0 {
			strerr := "tablero ya tiene slots generados"
			logError(strerr, err)
			return errors.New(strerr)
		}

		strSlots := ""
		for l := 0; l < 5; l++ {
			for r := 0; r < 5; r++ {
				str, err := board.generateSlot(l, r, !BOARD_HAS_CENTER && l == 2 && r == 2)
				if err != nil {
					strerr := fmt.Sprintf("failed generating slot %d %d", l, r)
					logError(strerr, err)
					return errors.New(strerr)
				}
				strSlots += str
			}
		}

		md = GetMD5Hash(strSlots)
		if game != nil && game.UniqueBoards {
			unique, _ = game.isUnique(string(md[:]))
		}
		if unique {
			break
		}
		if i == 99 {
			strerr := fmt.Sprintf("failed generating unique board for game")
			if game != nil {
				strerr += game.BingoID
			}
			logError(strerr, err)
			return errors.New(strerr)
		}
		err = board.deleteSlots()
		if err != nil {
			strerr := "failed to delete board slots"
			logError(strerr, nil)
			return errors.New(strerr)
		}
	}
	board.BoardHash = md
	if game.UniqueBoards && board.BoardHash == "" {
		strerr := fmt.Sprintf("empty hash on unique board mode: %s", game.BingoID)
		logError(strerr, nil)
		return errors.New(strerr)
	}
	err := board.guardar()
	if err != nil {
		strerr := fmt.Sprintf("failed saving board hash bid: %s", board.BoardID)
		logError(strerr, err)
		return errors.New(strerr)
	}
	return nil
}

// clearSlots desmarca todas las casillas del tablero
func (board *BingoBoard) clearSlots() error {
	var err = se.db.Table("bingo_slots").Where("bingo_id = ? AND board_id = ?", board.BingoID, board.BoardID).Update("marked", false).Error
	if !BOARD_HAS_CENTER {
		se.db.Table("bingo_slots").Where("bingo_id = ? AND board_id = ? AND letter = 'N' AND y = 2", board.BingoID, board.BoardID).Update("marked", true)
	}
	return err
}

// deleteSlots elimina las casillas
func (board *BingoBoard) deleteSlots() error {
	board.slots = []BingoSlot{}
	return se.db.Where("bingo_id = ? AND board_id = ?", board.BingoID, board.BoardID).Delete(&BingoSlot{}).Error
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

func (board *BingoBoard) drawImage(outName string) error {
	err := board.loadSlots()
	if err != nil {
		strerr := "error cargando slots @board.printText"
		logError(strerr, err)
		return errors.New(strerr)
	}

	const S = 2160
	im, err := gg.LoadImage("templateBoard.png")
	if err != nil {
		strerr := "failed opening template"
		logError(strerr, err)
		return errors.New(strerr)
	}

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
	dc.DrawImage(im, 0, 0)

	dc.SetRGB(255, 255, 255)
	if err := dc.LoadFontFace("moon_get-Heavy.ttf", 77); err != nil {
		panic(err)
	}

	dc.DrawStringAnchored(fmt.Sprint(board.ID), 215, 190, 0.5, 0.5)

	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("moon_get-Heavy.ttf", 96); err != nil {
		panic(err)
	}

	var (
		baseH float64 = 620
		baseW float64 = 560
		difH  float64 = 252
		difW  float64 = 252
	)

	rows := make(map[int][]int, 5)

	for _, slot := range board.slots {
		rows[letter2X(slot.Letter)] = append(rows[letter2X(slot.Letter)], slot.Number)
	}

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			if !(x == 2 && y == 2) || BOARD_HAS_CENTER {
				dc.DrawStringAnchored(strconv.Itoa(rows[x][y]), baseW+(difW*float64(x)), baseH+(difH*float64(y)), 0.5, 0.5)
			}
		}
	}

	outPath := "tableros/" + board.BingoID
	exists, err := pathExists(outPath)
	if err != nil {
		logError("failed checking directory "+outPath+" for drawn board", err)
	}
	if !exists {
		err := os.Mkdir(outPath, 0755)
		if err != nil {
			logError("failed creating directory "+outPath+" for drawn board", err)
		}
	}

	dc.Clip()
	dc.SavePNG(outPath + "/" + outName)

	fmt.Println("board file saved", outName)
	return nil
}

// markSlot marca una casilla (si la tiene)
// retorna si es ganador de la actual dinamica
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
			} else if dinamica == "n" {
				if lowerLetter == "b" || lowerLetter == "o" {
					haveInPlace++
				} else {
					x := letter2X(lowerLetter)
					if x == slot.Y {
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

// markSlot marca una casilla (si la tiene)
// retorna si es ganador de la actual dinamica
func (board *BingoBoard) markNCheck(drawn string, dinamica string) (bool, error) {
	var err error

	drawnList := strings.Split(drawn, ",")

	var winner bool
	for _, draw := range drawnList {
		winner, err = board.markSlots(string(draw[0]), string(draw[1:]), dinamica)
		if err != nil {
			strerr := fmt.Sprintf("failed checking board (%s) at game (%s)", board.BoardID, board.BingoID)
			logError(strerr, err)
			return winner, err
		}
	}

	return winner, nil
}
