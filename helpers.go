package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func logError(msg string, err error) {
	strerr := ""
	if err != nil {
		strerr = err.Error()
	}
	println(fmt.Sprintf("\033[31m Error %s: %s \033[39m", msg, strerr))
}

func getEnvOrDefault(llave string, defecto string) string {
	valor := os.Getenv(llave)
	if valor == "" {
		return defecto
	}
	return valor
}

// GetMD5Hash hashes a string using MD5.
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// GenerateRandomStringURLSafe returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

// UIDNew generates a URL safe string.
//  id := Gen(10)
//  fmt.Println(id)
//  // 9BZ1sApAX4
// Based on https://github.com/zemirco/uid/blob/master/uid.go
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// UIDNew takes constant letterBytes and returns random string of length n.
func UIDNew(n int) string {
	return UIDNewBytes(n, letters)
}

// UIDNewBytes takes letterBytes from parameters and returns random string of length n.
func UIDNewBytes(n int, lb string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(lb) {
			b[i] = lb[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func getGameNOrganizer(fromID, gameID string) (*BingoGame, *BingoOrganizer, error) {
	var (
		err       error
		game      BingoGame
		organizer BingoOrganizer
	)
	if gameID != "" {
		err = game.loadFromID(gameID)
		if err != nil {
			strerr := fmt.Sprintf("failed loading game (%s) from id for details", gameID)
			logError(strerr, err)
			return nil, nil, errors.New(strerr)
		}

		organizer, err = game.getOrganizer(fromID)
		if err != nil || organizer.ID == 0 && fromID != se.masterID {
			strerr := fmt.Sprintf("failed getting organizer TG:%s from game", fromID)
			logError(strerr, err)
			return nil, nil, errors.New(strerr)
		}
	} else {
		err = organizer.loadFromTG(fromID)
		if err != nil || organizer.ID == 0 {
			strerr := fmt.Sprintf("failed getting organizer TG:%s from tgid", fromID)
			logError(strerr, err)
			return nil, nil, errors.New(strerr)
		}

		game, err = organizer.getGame()

		if err != nil || organizer.ID == 0 {
			strerr := fmt.Sprintf("failed getting game from organizer %s", fromID)
			logError(strerr, err)
			return nil, nil, errors.New(strerr)
		}
	}
	return &game, &organizer, nil
}

// needFromDynamic retorna cuantas casillas deben estar marcadas
// ej: en 'linea' 5
// 	   en 'O' 16
func needFromDynamic(dinamica string) int {
	dinamica = strings.ToLower(dinamica)
	if len(dinamica) == 0 {
		logError("missing dinamica", nil)
		return 0
	}
	if dinamica == "a" {
		return 25
	} else if dinamica[0] == 'l' || dinamica[0] == '/' || dinamica[0] == '\\' {
		return 5
	} else if dinamica[0] == 'o' {
		return 16
	} else if dinamica[0] == 'u' {
		return 13
	} else if dinamica[0] == 'c' {
		return 13
	}
	return 0
}

func letter2X(b string) int {
	switch strings.ToLower(b) {
	case "b":
		return 0
	case "i":
		return 1
	case "n":
		return 2
	case "g":
		return 3
	case "o":
		return 4
	}
	return -1
}

func printSlot(n int) string {
	str := strconv.Itoa(n)
	if len(str) == 1 {
		str = " " + str
	}
	return str
}
