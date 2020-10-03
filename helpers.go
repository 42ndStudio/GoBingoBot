package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func logError(msg string, err error) {
	strerr := ""
	if err != nil {
		strerr = err.Error()
	}
	println(fmt.Sprintf("Error %s: %s", msg, strerr))
}

func getEnvOrDefault(llave string, defecto string) string {
	valor := os.Getenv(llave)
	if valor == "" {
		return defecto
	}
	return valor
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
		game      BingoGame
		organizer BingoOrganizer
	)
	err := game.loadFromID(gameID)
	if err != nil {
		strerr := fmt.Sprintf("failed loading game (%s) from id for details", gameID)
		logError(strerr, err)
		return nil, nil, errors.New(strerr)
	}
	if fromID != se.masterID {
		organizer, err = game.loadOrganizer(fromID)
		if err != nil || organizer.ID == 0 {
			strerr := fmt.Sprintf("failed getting organizer TG:%s from game", fromID)
			logError(strerr, err)
			return nil, nil, errors.New(strerr)
		}
	}

	return &game, &organizer, nil
}
