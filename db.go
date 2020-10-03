// Copyright 2020
// All Rights Reserved
// MuchLove
// 42nd Studio Team

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	//I&I justify it
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	authDBHost string
	authDBPort string
	authDBName string
	authDBUser string
	authDBPass string
)

func startDB() {
	var err error
	authDBHost = os.Getenv("DB_HOST")
	authDBPort = os.Getenv("DB_PORT")
	authDBName = os.Getenv("DB_NAME")
	authDBUser = os.Getenv("DB_USER")
	authDBPass = os.Getenv("DB_PASS")
	constring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", authDBUser, authDBPass, authDBHost, authDBPort, authDBName)
	se.db, err = gorm.Open("mysql", constring)
	if err != nil {
		fstring := fmt.Sprintf("failed to connect to database %s @ %s", authDBName, authDBHost)
		panic(fstring)
	}
	defer se.db.Close()
	autoMigrate()

	// Keep Connection Alive
	for {
		time.Sleep(time.Second * 350)
		var count int64
		se.db.Table("bingo_games").Where("created_at > NOW() - INTERVAL 24 HOUR").Count(&count)
		log.Println("games", zap.Int64("count", count))
	}
}

//autoMigrate corre las auto migraciones del servicio
// 	Gorm: AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func autoMigrate() {
	se.db.AutoMigrate(&BingoGame{})
	se.db.AutoMigrate(&BingoOrganizer{})
	se.db.AutoMigrate(&BingoBoard{})
	se.db.AutoMigrate(&BingoSlot{})
}
