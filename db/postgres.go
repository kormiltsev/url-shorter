package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var db *pg.DB

func StartConnection(conf [4]string) {
	db = pg.Connect(&pg.Options{
		Addr:     conf[0],
		User:     conf[1],
		Password: conf[2],
		Database: conf[3],
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
