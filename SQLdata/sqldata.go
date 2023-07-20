package sqldata

import (
	"database/sql"
	"log"
	gui "profitdetector/fynegui"

	_ "github.com/go-sql-driver/mysql"
)

func Getsqldb() {
	db, err := sql.Open("mysql", "user:1234@tcp(localhost:3306)/mycoinlist?charset=utf8")
	if err != nil {
		log.Println("Failed to insert data:", err)
		return
	}
	defer db.Close()

	log.Println("db structure complete")

	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect:", err)
		return
	}
	log.Println("db connect")

	//呼叫創建Fyne的涵式，並把database結構傳入
	gui.Createfyne(db)
}
