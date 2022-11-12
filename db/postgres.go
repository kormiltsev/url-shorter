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

func StartConnection() *pg.DB {
	godotenv.Load()
	adr := os.Getenv("ADR")
	usr := os.Getenv("USR")
	pwd := os.Getenv("PASS")
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
	// check connection
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		panic(err)
	}
	// ==============
	return db
	//defer db.Close()
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func RowsQuantity(db *pg.DB, tablename string) (int, error) {
	var qty int
	_, err := db.Query(&qty, fmt.Sprintf(`SELECT COUNT(surl) FROM %s`, tablename))
	return qty, err
}

func ReturnTablesNames(db *pg.DB) ([]string, error) {
	var qry []string
	_, err := db.Query(&qry, fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_schema NOT IN ('information_schema','pg_catalog')`))
	return qry, err
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

func InsertNewRow(db *pg.DB, newrow *Baserow) {
	err := db.Insert(newrow)
	if err != nil {
		panic(err)
	}
}

func CreateTable(db *pg.DB) {
	var urler Baserow
	err := db.CreateTable(&urler, &orm.CreateTableOptions{
		Temp:          false,
		IfNotExists:   true,
		FKConstraints: true,
	})
	panicIf(err)
}
