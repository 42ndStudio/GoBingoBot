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
	authDBHost = os.Getenv("db_host")
	authDBPort = os.Getenv("db_port")
	authDBName = os.Getenv("db_name")
	authDBUser = os.Getenv("db_user")
	authDBPass = os.Getenv("db_pass")
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
		se.db.Table("system_events").Where("created_at > NOW() - INTERVAL 24 HOUR").Count(&count)
		log.Println("last_24h_system_events", zap.Int64("count", count))
	}
}

//autoMigrate corre las auto migraciones del servicio
// 	Gorm: AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func autoMigrate() {
	se.db.AutoMigrate(&BingoGame{})
	se.db.AutoMigrate(&BingoClient{})
	se.db.AutoMigrate(&BingoBoard{})
	se.db.AutoMigrate(&BingoSlot{})
}
