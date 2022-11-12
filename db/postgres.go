package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/joho/godotenv"
	// "github.com/go-pg/pg"
	//"github.com/go-pg/pg/orm"
	// "github.com/joho/godotenv"
)

var db *pg.DB

func StartConnection() error {
	godotenv.Load()
	adr := os.Getenv("ADR")
	usr := os.Getenv("USR")
	pwd := os.Getenv("PASS")
	dbs := os.Getenv("DB")

	if adr == "" || usr == "" || pwd == "" || dbs == "" {
		log.Printf("Can't find DB specs in .env")
	}
	db = pg.Connect(&pg.Options{
		Addr:     adr,
		User:     usr,
		Password: pwd,
		Database: dbs,
	})
	// check connection
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		panic(err)
	}
	// ==============
	return err
	//defer db.Close()
}

func CreateTable() error {
	var urler Baserow
	err := db.CreateTable(&urler, &orm.CreateTableOptions{
		Temp:          false,
		IfNotExists:   true,
		FKConstraints: true,
	})
	panicIf(err)
	return err
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

func ReturnTablesInfo(db *pg.DB, tablename string) {
	var info []struct {
		ColumnName string
		DataType   string
	}
	_, err := db.Query(&info, fmt.Sprintf(`
	SELECT column_name, data_type
	FROM information_schema.columns
	WHERE table_name = surl`)) //, tablename))
	panicIf(err)
	fmt.Println(info)
}

func InsertNewRow(newrow *Baserow) {
	err := db.Insert(newrow)
	if err != nil {
		panic(err)
	}
}

func QueryGetURL(row *Baserow) error {
	//var row Baserow
	err := db.Model(row).
		ColumnExpr("lurl").
		Where("surl = ?", row.Surl).
		Select()

		// to check:
	// item := Model2{
	// 	Id: 2,
	// }
	// err := db.Select(row)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(item)
	return err
}
