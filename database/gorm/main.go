package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Album struct {
	Id     int    `gorm:"primarykey"`
	Title  string `gorm:"size=128"`
	Artist string `gorm:"size=255"`
	Price  float32
}

func main() {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(127.0.0.1:3306)/recordings?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("DBUSER"),
			os.Getenv("DBPASS"),
		), // data source name
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")

	var album = Album{Id: 1}

	fmt.Println("Search ID=1")
	db.First(&album).Scan(&album)
	fmt.Println(album)

	var albums []Album

	fmt.Println("Search All")
	db.Find(&albums).Scan(&albums)
	for _, v := range albums {
		fmt.Println(v)
	}
	var newalbum = Album{
		Title:  "Song sing",
		Price:  12.31,
		Artist: "Joana",
	}

	fmt.Println("Create new data")
	db.Create(&newalbum)

	db.Find(&albums).Scan(&albums)
	for _, v := range albums {
		fmt.Println(v)
	}

	fmt.Println("Search by title")
	db.Where(&Album{Title: "Blue Train"}).Find(&albums)

	for _, v := range albums {
		fmt.Println(v)
	}

}
