// Copyright 2020
// All Rights Reserved
// MuchLove
// 42nd Studio Team

package main

import (
	"github.com/jinzhu/gorm"

	//I&I justify it
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// BingoGame es un evento de bingo el cual tiene muchos tableros asociados
type BingoGame struct {
	gorm.Model
	BingoID             string `gorm:"unique;not null"`
	Name                string `gorm:"not null"`
	CurrentMode         string `gorm:"null"`
	BoardsSold          int64  `gorm:"not null" sql:"DEFAULT:0"`
	Password            string `gorm:"not null"`
	AcceptingOrganizers bool   `gorm:"not null"`
}

// BingoClient es un evento de bingo el cual tiene muchos tableros asociados
type BingoClient struct {
	gorm.Model
	BingoID    string `gorm:"unique;not null"` // A que bingo tiene acceso
	Name       string `gorm:"not null"`
	TelegramID string `gorm:"not null"`
	BoardsSold int    `gorm:"null"`
}

// BingoBoard es un evento de bingo el cual tiene muchos tableros asociados
type BingoBoard struct {
	gorm.Model
	BingoID     string `gorm:"not null"`
	BoardID     string `gorm:"not null"`
	Name        string `gorm:"null"`
	GamesWon    int    `gorm:"null"`
	GamesWonIds string `gorm:"null"`
	Sold        bool   `gorm:"null"`
}

// BingoSlot es una casilla de un tablero
type BingoSlot struct {
	gorm.Model
	BoardID string `gorm:"not null"`
	Letter  string `gorm:"not null"`
	Number  int    `gorm:"not null"`
	Marked  bool   `gorm:"not null" sql:"DEFAULT:'false'"`
}
