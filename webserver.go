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

type RequestBoardGenerate struct {
	BingoID string
	Boards  int
}

type RequestBoardCheck struct {
	BingoID string
	BoardId string
}

type RequestDrawBalot struct {
	BingoID string
	Balot   string
}

type RequestSetMode struct {
	BingoID string
	Param   string
}

type ResponseGame struct {
	BingoID     string
	Name        string
	Playing     bool
	CurrentMode string
	BoardsSold  int
	DrawnBalots string
}

type ResponseBoardCheck struct {
	Winner      bool
	CurrentMode string
	DrawnBalots string
}

type ResponseStd struct {
	Status  string
	Message string
}

// dataresp := []*responseContentData{}
// for _, valor := range valores.Datas {
// 	dataresp = append(dataresp, &responseContentData{
// 		ID:         valor.Id,
// 		ObjectID:   valor.ObjectId,
// 		Value:      valor.Value,
// 		Identifier: valor.Identifier,
// 	})
// }

// // Escribir respuesta
// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusOK)
// json.NewEncoder(w).Encode(dataresp)

func handleGames(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		game BingoGame
	)
	mutex.Lock()

	if r.Method == "GET" {
		gid := r.URL.Query().Get("gid")

		if gid == "" {
			// Load All Games
			var (
				games []BingoGame
				rows  []ResponseGame
			)
			err = se.db.Where("1 = 1").Find(&games).Error
			if err != nil {
				logError("failed loading games @gamesList", nil)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
			for _, game := range games {
				rows = append(rows, ResponseGame{
					Name:        game.Name,
					BingoID:     game.BingoID,
					Playing:     game.Playing,
					CurrentMode: game.CurrentMode,
					DrawnBalots: game.DrawnBalots,
					BoardsSold:  game.BoardsSold,
				})
			}
			// Escribir respuesta
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(rows)
		} else {
			// Load Specific Game
			err = game.loadFromID(gid)
			if err != nil {
				strerr := fmt.Sprintf("failed loading game (%s) from id @handleGames[GET]", gid)
				fmt.Println(err)
				fmt.Println(gorm.ErrRecordNotFound)
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}

			resp := ResponseGame{
				BingoID:     game.BingoID,
				Playing:     game.Playing,
				Name:        game.Name,
				BoardsSold:  game.BoardsSold,
				DrawnBalots: game.DrawnBalots,
				CurrentMode: game.CurrentMode,
			}

			// Escribir respuesta
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}

	} else if r.Method == "POST" {
		var (
			inputGame BingoGame
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
				strerr := fmt.Sprintf("failed loading game (%s) from id @handleGames[POST]", inputGame.BingoID)
				fmt.Println(err)
				fmt.Println(gorm.ErrRecordNotFound)
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
		}

		game.Name = inputGame.Name
		game.IdentifierType = inputGame.IdentifierType

		fmt.Println("saving game")
		err = game.guardar()
		if err != nil {
			strerr := "failed game.guardar()"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
		}

		// Escribir respuesta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(game)
	}

	mutex.Unlock()
}

func handleBoardGenerate(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "POST" {
		var (
			err      error
			reqInput RequestBoardGenerate
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
			strerr := fmt.Sprintf("failed loading game (%s) from id @handleBoardGenerate", reqInput.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Generate board
		board, err := game.generateBoard()
		go func() {
			fmt.Println("drawing board")
			err = board.drawImage(board.BoardHash + ".png")
			if err != nil {
				strerr := fmt.Sprintf("failed drawing board (%s) for game (%s)", board.BoardID, reqInput.BingoID)
				logError(strerr, err)
				return
			}
		}()

		// Escribir respuesta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(board)
	}

	mutex.Unlock()
}

func handleGameSetMode(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "POST" {
		var (
			err      error
			reqInput RequestSetMode
			game     BingoGame
			resp     ResponseStd
		)

		// Obtain Input Values
		err = json.NewDecoder(r.Body).Decode(&reqInput)
		if err != nil {
			strerr := "failed decoding input @handleGameSetMode"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Load Bingo Game
		err = game.loadFromID(reqInput.BingoID)
		if err != nil {
			strerr := fmt.Sprintf("failed loading game (%s) from id @handleGameSetMode", reqInput.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		fmt.Println("handleGameSetMode", game.BingoID, reqInput.Param)

		// Set Mode
		// Its a different game mode
		if reqInput.Param == "PLAY" || reqInput.Param == "STOP" {
			game.Playing = reqInput.Param == "PLAY"
			err = game.guardar()
			if err != nil {
				strerr := fmt.Sprintf("failed saving game (%s) as playing = %v  @handleGameSetMode", reqInput.BingoID, reqInput.Param == "PLAY")
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
			fmt.Println("game ", game.BingoID, "saved as playing = ", reqInput.Param == "PLAY")
			msg := "Playing"
			if !game.Playing {
				msg = "Not Playing"
			}
			resp = ResponseStd{
				Status:  "OK",
				Message: msg,
			}
		} else if reqInput.Param == "CLEAR_BALOTS" {
			ogdrawn, err := game.clearSlots()
			if err != nil {
				strerr := fmt.Sprintf("failed clearing game (%s) slots  @handleGameSetMode", reqInput.BingoID)
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
			fmt.Println("cleared balots for game", game.BingoID, "drawn:", game.DrawnBalots)
			resp = ResponseStd{
				Status:  "OK",
				Message: ogdrawn,
			}
		} else if stringInSlice(reqInput.Param, GAME_MODES) {
			game.CurrentMode = reqInput.Param
			err = game.guardar()
			if err != nil {
				strerr := "failed saving game new mode @handleGameSetMode"
				logError(strerr, err)
				http.Error(w, RESP_ERROR, http.StatusBadRequest)
				mutex.Unlock()
				return
			}
			fmt.Println("changed game", game.BingoID, "to mode:", game.CurrentMode)
			resp = ResponseStd{
				Status:  "OK",
				Message: game.CurrentMode,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
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
		_, err = game.drawBalot(string(reqInput.Balot[0]), string(reqInput.Balot[1:]))
		if err != nil {
			strerr := fmt.Sprintf("failed game.drawBalot reqInput.Balot of bingo %s", game.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		resp := ResponseStd{
			Status:  "OK",
			Message: reqInput.Balot,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	}

	mutex.Unlock()
}

func handleBoardCheck(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()

	if r.Method == "POST" {
		var (
			err      error
			reqInput RequestBoardCheck
			game     BingoGame
			resp     ResponseBoardCheck
		)

		// Obtain Input Values
		err = json.NewDecoder(r.Body).Decode(&reqInput)
		if err != nil {
			strerr := "failed decoding input @handleBoardCheck"
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Load Bingo Game
		err = game.loadFromID(reqInput.BingoID)
		if err != nil {
			strerr := fmt.Sprintf("failed loading game (%s) from id @handleBoardCheck", reqInput.BingoID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Load board
		board, err := game.getBoard(reqInput.BoardId)
		if err != nil {
			strerr := fmt.Sprintf("failed loading board. game: %s board: %s @handleBoardCheck", reqInput.BingoID, reqInput.BoardId)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		// Mark drawn slots and check if winner
		resp.Winner, err = board.markNCheck(game.DrawnBalots, game.CurrentMode)
		if err != nil {
			strerr := fmt.Sprintf("failed board %s markNCheck", board.BoardID)
			logError(strerr, err)
			http.Error(w, RESP_ERROR, http.StatusBadRequest)
			mutex.Unlock()
			return
		}

		resp.DrawnBalots = game.DrawnBalots
		resp.CurrentMode = game.CurrentMode

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

	mutex.Unlock()
}

func runWebServer() {
	http.HandleFunc("/games/", handleGames)

	http.HandleFunc("/games/setmode/", handleGameSetMode)

	http.HandleFunc("/games/drawbalot/", handleDrawBalot)

	http.HandleFunc("/games/boards/generate/", handleBoardGenerate)

	http.HandleFunc("/games/boards/check/", handleBoardCheck)

	fmt.Println("starting web server @", WEB_SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+WEB_SERVER_PORT, nil))
}
