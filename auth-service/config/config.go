package config

import (
	"auth-service/routes"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	//TODO driver mysql
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Db *sql.DB
}

func InitConfig() *Config {
	db := initDatabase()
	return &Config{
		Db: db,
	}
}

func (app *Config) InitServer() {

	server := &http.Server{
		Addr:    ":80",
		Handler: routes.Routes(app.Db),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func initDatabase() *sql.DB {
	var db *sql.DB
	var counts int64

	conn, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		err = conn.Ping()
		if err != nil {
			log.Println(err)
			counts++

			if counts > 10 {
				break
			}

			time.Sleep(5 * time.Second)
			continue
		} else {
			db = conn
			break
		}
	}

	log.Println("connection established")
	return db
}
