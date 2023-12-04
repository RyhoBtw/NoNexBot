package database

import (
	"NoiseDcBot"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDatabase() *sql.DB {
	log.Println(dsn())
	db := OpenDB()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println(pingErr)
	}

	defer db.Close()

	return db
}

func dsn() string {
	conf := NoiseDcBot.GetConf()
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", conf.DBUser, conf.DBPassword, conf.Host, conf.DBName)
}

func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	return db
}
