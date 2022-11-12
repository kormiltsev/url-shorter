package db

import (
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
	// "github.com/go-pg/pg"
	// "github.com/joho/godotenv"
)

func StartConnection() *pg.DB {
	godotenv.Load()
	adr := os.Getenv("ADR")
	usr := os.Getenv("USER")
	pwd := os.Getenv("PWD")
	dbs := os.Getenv("DB")

	if adr == "" || usr == "" || pwd == "" || dbs == "" {
		log.Printf("Can't find DB specs in .env")
	}
	db := pg.Connect(&pg.Options{
		Addr:     adr,
		User:     usr,
		Password: pwd,
		Database: dbs,
	})
	return db
	//defer db.Close()
}

func CalcQty(db *pg.DB, tablename string) (int, error) {
	var qty int
	_, err := db.Query(&qty, fmt.Sprintf(`SELECT COUNT(data_source) FROM %s`, tablename))
	return qty, err

}
