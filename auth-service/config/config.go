package config

import (
	"auth-service/routes"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int

type Config struct {
	Db *sql.DB
}

func InitConfig() *Config {
	db := connectToDB()
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

func _() *sql.DB {
	var db *sql.DB
	var counts int64

	conn, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("postgres no parece estar listo...")
			counts++
		} else {
			log.Println("conectado a postgress")

			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("esperando por dos segundos")
		time.Sleep(2 * time.Second)
		continue
	}
}
