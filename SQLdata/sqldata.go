package sqldata

import (
	"database/sql"
	"fmt"
	gui "profitdetector/fynegui"

	_ "github.com/go-sql-driver/mysql"
)

type coins struct {
	name    string
	amount  int
	cost    int
	avgcost int
}

func Getsqldb() {
	db, err := sql.Open("mysql", "user:1234@tcp(localhost:3306)/mycoinlist?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("db structure complete")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("db conneected")

	//呼叫創建Fyne的涵式，並把database結構傳入
	gui.Createfyne(db)
}
