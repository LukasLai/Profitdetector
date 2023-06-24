package fynegui

import (
	"database/sql"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	_ "github.com/go-sql-driver/mysql"
)

type Total struct {
	Totalamount   int
	Totalcost     int
	Totalavgprice float64
}

type Inputdata struct {
	Name   string
	Cost   float64
	Amount float64
	Avg    float64
}

var curPriceText *canvas.Text

func Createfyne(db *sql.DB) {
	//把總量、總花費以及平均儲存進去結構
	t := Total{} //儲存總資料的結構
	err := db.QueryRow("SELECT SUM(amount) FROM cryptolist").Scan(&t.Totalamount)
	if err != nil {
		panic(err)
	}

	err = db.QueryRow("SELECT SUM(cost) FROM cryptolist").Scan(&t.Totalcost)
	if err != nil {
		panic(err)
	}
	err = db.QueryRow("SELECT SUM(cost)/SUM(amount) FROM cryptolist").Scan(&t.Totalavgprice)
	if err != nil {
		panic(err)
	}

	//=========================================================================================================
	a := app.New()
	w := a.NewWindow("Window Example")
	w.Resize(fyne.NewSize(500, 500))

	// 上方視窗;
	inputdata := Inputdata{} //儲存輸入資料的結構

	//大標題
	topic := canvas.NewText("Profid Detector", color.NRGBA{R: 120, G: 0, B: 120, A: 255})
	topic.Alignment = fyne.TextAlignCenter
	topic.TextSize = 20

	topic1 := canvas.NewText("Current Update", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	topic1.Alignment = fyne.TextAlignLeading // 將對齊方式設定為置左
	topic1.TextSize = 20

	label := canvas.NewText("Name: BTC        Amount: "+fmt.Sprintf("%.2f", inputdata.Amount)+
		"        Cost: $"+fmt.Sprintf("%f", inputdata.Cost)+
		"        Averageprice: "+fmt.Sprintf("%.2f", inputdata.Avg), color.Black)
	label.Alignment = fyne.TextAlignCenter
	label.TextSize = 15
	//現價顯示
	curPriceText = canvas.NewText("Current Price: $0.00", color.Black)
	curPriceText.Alignment = fyne.TextAlignLeading
	curPriceText.TextSize = 15

	topWindow := container.NewVBox(topic, topic1, label, curPriceText)

	// 中間視窗
	nameEntry := widget.NewEntry()
	nameLabel := canvas.NewText("Name:", color.Black)
	nameLabel.TextSize = 15
	nameContainer := container.NewHBox(nameLabel, nameEntry)

	costEntry := widget.NewEntry()
	costLabel := canvas.NewText("Cost:", color.Black)
	costLabel.TextSize = 15
	costContainer := container.NewHBox(costLabel, costEntry)

	amountEntry := widget.NewEntry()
	amountLabel := canvas.NewText("Amount:", color.Black)
	amountLabel.TextSize = 15
	amountContainer := container.NewHBox(amountLabel, amountEntry)

	updateLabel := canvas.NewText("", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	updateLabel.Alignment = fyne.TextAlignCenter
	updateLabel.TextSize = 15
	updateContainer := container.NewHBox(updateLabel)
	//送出按鈕，把資料輸入至資料庫儲存
	sentButton := widget.NewButton("Sent", func() {
		inputdata.Name = nameEntry.Text
		cost := costEntry.Text
		amount := amountEntry.Text
		floatAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			fmt.Println("無法轉換為 float64")
			return
		}
		floatCost, err := strconv.ParseFloat(cost, 64)
		if err != nil {
			fmt.Println("無法轉換為 float64")
			return
		}
		avg := floatCost / floatAmount

		inputdata.Amount = floatAmount
		inputdata.Cost = floatCost
		inputdata.Avg = avg
		label.Text = "Name: " + inputdata.Name + "        Amount: " + fmt.Sprintf("%.2f", inputdata.Amount) +
			"        Cost: $" + fmt.Sprintf("%.2f", inputdata.Cost) +
			"        Averageprice: " + fmt.Sprintf("%.2f", inputdata.Avg)
		insertinfo(db, inputdata.Name, inputdata.Amount, inputdata.Cost, inputdata.Avg)
		updateLabel.Text = "Data Updated!"
		updateLabel.Refresh()
		label.Refresh()
	})

	middleWindow := container.NewHBox(nameContainer, costContainer, amountContainer, sentButton, updateContainer)

	// 下方視窗; 負責顯示所有總量的資料
	curprice := canvas.NewText("", color.NRGBA{R: 0, G: 120, B: 120, A: 255})
	curprice.Alignment = fyne.TextAlignTrailing
	curprice.TextSize = 15

	topic2 := canvas.NewText("TotalAsset", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	topic2.Alignment = fyne.TextAlignLeading
	topic2.TextSize = 20

	bottomLabel := canvas.NewText("Name: BTC       Amount: "+strconv.Itoa(t.Totalamount)+
		"       Cost: $"+strconv.Itoa(t.Totalcost)+"        Avgprice :"+fmt.Sprintf("%.2f", t.Totalavgprice), color.Black)
	bottomLabel.Alignment = fyne.TextAlignCenter
	bottomLabel.TextSize = 15
	bottomWindow := container.NewVBox(curprice, topic2, bottomLabel)

	// 將視窗組合並排列
	content := container.NewVBox(
		container.NewVBox(topWindow),
		layout.NewSpacer(),
		container.NewVBox(middleWindow),
		layout.NewSpacer(),
		container.NewVBox(bottomWindow),
	)

	w.SetContent(content)
	w.ShowAndRun()
}

/*func Currentprice(price string) {
	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("無法轉換為 float64")
		return
	}

}*/

// 把此次新增的交易新增到database
func insertinfo(db *sql.DB, name string, amount, cost, avgcost float64) {

	insertStmt, err := db.Prepare("INSERT INTO cryptolist(name,amount,cost,avgcost) VALUES(?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = insertStmt.Exec(name, amount, cost, avgcost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert complete")
}

func Updatecurrentprice(price string) {
	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println("無法轉換為 float64")
		return
	}
	fmt.Println("順利進入Updatacurrentprice")
	curPriceText.Text = "Current Price: $" + fmt.Sprintf("%.2f", floatPrice)
	curPriceText.Refresh()
}
