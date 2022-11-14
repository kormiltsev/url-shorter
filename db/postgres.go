package db

import (
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/joho/godotenv"
)

var db *pg.DB

func StartConnection() {
	godotenv.Load()
	adr := os.Getenv("ADR")
	usr := os.Getenv("USR")
	pwd := os.Getenv("PASS")
	dbs := os.Getenv("DB")

	if adr == "" || usr == "" || pwd == "" || dbs == "" {
		adr = "localhost:5432"
		usr = "postgres"
		pwd = "root"
		dbs = "postgres"
	}
	db = pg.Connect(&pg.Options{
		Addr:     adr,
		User:     usr,
		Password: pwd,
		Database: dbs,
	})
	CreateTable()
}

func DbClose() {
	db.Close()
}

func CreateTable() {
	var urler Baserow
	err := db.CreateTable(&urler, &orm.CreateTableOptions{
		Temp:          false,
		IfNotExists:   true,
		FKConstraints: true,
	})
	panicIf(err)
	//return err
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func RowsQuantity() (int, error) {
	var qty int
	_, err := db.Query(&qty, `SELECT COUNT(surl) FROM baserows`)
	return qty, err
}

func QueryGetURL(row *Baserow) error {
	//var row Baserow
	return db.Model(row).
		ColumnExpr("lurl").
		Where("surl = ?", row.Surl).
		Select()
}

func QueryPOSTorSelect(row *Baserow) error {
	_, err := db.Model(row). // true or false
					ColumnExpr("surl").
					Where("lurl = ?", row.Lurl).
					OnConflict("DO NOTHING"). // optional
					SelectOrInsert()
	return err
}
