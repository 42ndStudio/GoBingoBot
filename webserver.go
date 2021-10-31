// With love
// 42nd Studio
// 2020-2021

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"gorm.io/gorm"
)

var mutex = &sync.Mutex{}

const WEB_SERVER_PORT = "8042"
const RESP_ERROR = "An error occurred\nOcurrió un error"

type RequestGenerateBoard struct {
	BingoID string
	Boards  int
}

type RequestDrawBalot struct {
	BingoID string
	Balot   string
}

func handleGames(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "GET" {
		gid := r.URL.Query().Get("gid")
		fmt.Fprintf(w, "its a GET gid: "+gid)
	} else if r.Method == "POST" {
		var (
			err       error
			inputGame BingoGame
			game      BingoGame
		)

		// Obtain Input Values
		err = json.NewDecoder(r.Body).Decode(&inputGame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Check for saved game
		//  if not exists create
		err = game.loadFromID(inputGame.BingoID)
		if err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				fmt.Println("creating game")
			} else {
				strerr := fmt.Sprintf("failed loading game (%s) from id", inputGame.BingoID)
				fmt.Println(err)
				fmt.Println(gorm.ErrRecordNotFound)
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
		}

		game.Name = inputGame.Name
		fmt.Println("saving game")
		err = game.guardar()
		if err != nil {
			strerr := "failed game.guardar()"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
		}

		fmt.Fprintf(w, "its a POST game: "+game.Name)
		// fmt.Fprintf(w, fmt.Sprintf("%#v", game))

		// game := new(BingoGame)
		// err := game.guardar()
		// if err != nil {
		// 	strerr := "failed game.guardar()"
		// 	logError(strerr, err)
		// } else {
		// 	organizer := new(BingoOrganizer)
		// 	organizer.TelegramID = fromID
		// 	organizer.BingoID = game.BingoID

		// 	err = organizer.guardar()
		// 	if err != nil {
		// 		strerr := "failed organizer.guardar()"
		// 		logError(strerr, err)
		// 	} else {
		// 		respmsg.Text = fmt.Sprintf("Juego %s creado", game.BingoID)
		// 		respmsg.ReplyMarkup = masterKeyboard
		// 		delete(waitingon, fromID)
		// 	}
		// }
	}

	mutex.Unlock()
}

func handleBoardGenerate(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "POST" {
		var (
			err      error
			reqInput RequestGenerateBoard
			game     BingoGame
		)

		// Obtain Input Values
		err = json.NewDecoder(r.Body).Decode(&reqInput)
		if err != nil {
			strerr := "failed decoding input @handleBoardGenerate"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Load Bingo Game
		err = game.loadFromID(reqInput.BingoID)
		if err != nil {
			strerr := fmt.Sprintf("failed loading game (%s) from id", reqInput.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Generate board
		board, err := game.generateBoard()
		go func() {
			fmt.Println("generating board")
			err = board.drawImage(board.BoardHash + ".png")
			if err != nil {
				strerr := fmt.Sprintf("failed drawing board (%s) for game (%s)", board.BoardID, reqInput.BingoID)
				logError(strerr, err)
				return
			}
		}()
	}

	mutex.Unlock()
}

func handleDrawBalot(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "POST" {
		var (
			err      error
			reqInput RequestDrawBalot
			game     BingoGame
		)

		// Obtain Input Values
		err = json.NewDecoder(r.Body).Decode(&reqInput)
		if err != nil {
			strerr := "failed decoding input @handleBoardGenerate"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Load Bingo Game
		err = game.loadFromID(reqInput.BingoID)
		if err != nil {
			strerr := fmt.Sprintf("failed loading game (%s) from id", reqInput.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Draw balot
		_, err = game.drawBalot(string(reqInput.Balot[0]), string(reqInput.Balot[1]))
		if err != nil {
			strerr := fmt.Sprintf("failed game.drawBalot reqInput.Balot of bingo %s", game.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

	}

	mutex.Unlock()
}

func runWebServer() {
	http.HandleFunc("/games/", handleGames)

	http.HandleFunc("/games/boards/generate/", handleBoardGenerate)

	http.HandleFunc("/games/drawbalot/", handleDrawBalot)

	log.Fatal(http.ListenAndServe(":"+WEB_SERVER_PORT, nil))
}
