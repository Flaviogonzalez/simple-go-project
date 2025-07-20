package config

import (
	"database/sql"
	"log"
	"log-service/routes"
	"net/http"
	"os"
	"time"
)

var counts int
var port = ":80"

type Config struct {
	db *sql.DB
}

func InitConfig() *Config {
	conn := connectToDB()

	return &Config{
		db: conn,
	}
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

func (app *Config) InitServer() {
	srv := &http.Server{
		Addr:    port,
		Handler: routes.Routes(app.db),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
