// Copyright 2020
// All Rights Reserved
// MuchLove
// 42nd Studio Team

package main

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	authDBHost string
	authDBPort string
	authDBName string
	authDBUser string
	authDBPass string
)

// Product -- Represents a product
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// TableName setting the table name
func (Product) TableName() string {
	return "allProducts"
}

func startSQLLiteDB() {
	var err error
	se.db, err = gorm.Open("sqlite3", "bingo.db")
	if err != nil {
		panic("failed to connect database")
	}
	// defer db.Close()
	// var product Product
	// rows, err := db.Model(&Product{}).Rows()
	// defer rows.Close()
	// if err != nil {
	//     panic(err)
	// }
	// for rows.Next() {
	//     db.ScanRows(rows, &product)
	//     fmt.Println(product)
	// }

	autoMigrate()
}

// func startSQLLiteDB() {
// 	var err error

// 	// github.com/mattn/go-sqlite3
// 	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

// 	if err != nil {
// 		fstring := fmt.Sprintf("failed to connect to database %s @ %s", authDBName, authDBHost)
// 		panic(fstring)
// 	}
// 	defer se.db.Close()

// 	autoMigrate()

// }

// func startDBFunk() {
// 	db, err := sql.Open("sqlite3", "./foo.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	log.Println("have foo?")

// 	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='foo';")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	have := 0
// 	for rows.Next() {
// 		var name string
// 		err = rows.Scan(&name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(name)
// 		have += 1
// 	}

// 	sqlStmt := `
// 	create table foo (id integer not null primary key, name text);
// 	delete from foo;
// 	`
// 	_, err = db.Exec(sqlStmt)
// 	if err != nil {
// 		log.Printf("%q: %s\n", err, sqlStmt)
// 		return
// 	}

// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	for i := 0; i < 100; i++ {
// 		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	tx.Commit()

// 	rows, err = db.Query("select id, name from foo")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err = rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, name)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	stmt, err = db.Prepare("select name from foo where id = ?")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
// 	var name string
// 	err = stmt.QueryRow("3").Scan(&name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(name)

// 	_, err = db.Exec("delete from foo")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	rows, err = db.Query("select id, name from foo")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		err = rows.Scan(&id, &name)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(id, name)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func startDB() {
// 	var err error
// 	authDBHost = os.Getenv("DB_HOST")
// 	authDBPort = os.Getenv("DB_PORT")
// 	authDBName = os.Getenv("DB_NAME")
// 	authDBUser = os.Getenv("DB_USER")
// 	authDBPass = os.Getenv("DB_PASS")
// 	constring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", authDBUser, authDBPass, authDBHost, authDBPort, authDBName)
// 	se.db, err = gorm.Open("mysql", constring)
// 	if err != nil {
// 		fstring := fmt.Sprintf("failed to connect to database %s @ %s", authDBName, authDBHost)
// 		panic(fstring)
// 	}
// 	defer se.db.Close()
// 	autoMigrate()

// 	// Keep Connection Alive
// 	for {
// 		time.Sleep(time.Second * 350)
// 		var count int64
// 		se.db.Table("bingo_games").Where("created_at > NOW() - INTERVAL 24 HOUR").Count(&count)
// 		log.Println("games", zap.Int64("count", count))
// 	}
// }

//autoMigrate corre las auto migraciones del servicio
// 	Gorm: AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
func autoMigrate() {

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	se.db.AutoMigrate(&BingoGame{})
	se.db.AutoMigrate(&BingoOrganizer{})
	se.db.AutoMigrate(&BingoBoard{})
	se.db.AutoMigrate(&BingoSlot{})
}
