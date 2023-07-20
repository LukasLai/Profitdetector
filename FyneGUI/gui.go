package fynegui

import (
	"database/sql"
	"fmt"
	"image/color"
	"log"
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
	Totalamount   float64
	Totalcost     float64
	Totalavgprice float64
}

type Inputdata struct {
	Name   string
	Cost   float64
	Amount float64
	Avg    float64
}

var curPriceText *canvas.Text
var curValue *canvas.Text
var curProfit *canvas.Text
var totalAmountProift float64
var totalCost float64

func Createfyne(db *sql.DB) {
	//把總量、總花費以及平均儲存進去結構
	t := Total{} //儲存總資料的結構
	err := db.QueryRow("SELECT SUM(amount) FROM cryptolist").Scan(&t.Totalamount)
	if err != nil {
		log.Println("Failed to Select SUM(amount): ", err)
		return
	}
	totalAmountProift = t.Totalamount

	err = db.QueryRow("SELECT SUM(cost) FROM cryptolist").Scan(&t.Totalcost)
	if err != nil {
		log.Println("Failed to Select SUM(cost): ", err)
		return
	}
	totalCost = t.Totalcost
	err = db.QueryRow("SELECT SUM(cost)/SUM(amount) FROM cryptolist").Scan(&t.Totalavgprice)
	if err != nil {
		log.Println("Failed to Select SUM(cost)/SUM(amount): ", err)
		return
	}

	//=========================================================================================================
	a := app.New()
	w := a.NewWindow("BTC Profitdector")
	w.Resize(fyne.NewSize(700, 500))

	// 上方視窗;
	inputdata := Inputdata{} //儲存輸入資料的結構

	//大標題
	topic := canvas.NewText("Profid Detector", color.NRGBA{R: 120, G: 0, B: 120, A: 255})
	topic.Alignment = fyne.TextAlignCenter
	topic.TextSize = 20

	topic1 := canvas.NewText("Current Update", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	topic1.Alignment = fyne.TextAlignLeading // 將對齊方式設定為置左
	topic1.TextSize = 20

	label := canvas.NewText("Name: BTC        Amount: "+fmt.Sprintf("%.4f", inputdata.Amount)+
		"        Cost: $"+fmt.Sprintf("%f", inputdata.Cost)+
		"        Averageprice: "+fmt.Sprintf("%.4f", inputdata.Avg), color.Black)
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
	updateLabel.Alignment = fyne.TextAlignLeading
	updateLabel.TextSize = 15

	//送出按鈕，把資料輸入至資料庫儲存
	sentButton := widget.NewButton("Sent", func() {
		inputdata.Name = nameEntry.Text
		cost := costEntry.Text
		amount := amountEntry.Text
		floatAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Println("Failed to parse floatAmount: ", err)
			return
		}
		floatCost, err := strconv.ParseFloat(cost, 64)
		if err != nil {
			log.Println("Failed to parse floatCost: ", err)
			return
		}
		avg := floatCost / floatAmount

		inputdata.Amount = floatAmount
		inputdata.Cost = floatCost
		inputdata.Avg = avg
		label.Text = "Name: " + inputdata.Name + "        Amount: " + fmt.Sprintf("%.4f", inputdata.Amount) +
			"        Cost: $" + fmt.Sprintf("%.4f", inputdata.Cost) +
			"        Averageprice: " + fmt.Sprintf("%.4f", inputdata.Avg)
		insertinfo(db, inputdata.Name, inputdata.Amount, inputdata.Cost, inputdata.Avg)
		updateLabel.Text = "Data Updated!"
		updateLabel.Refresh()
		label.Refresh()
	})
	middleHBox := container.NewHBox(nameContainer, costContainer, amountContainer, sentButton)
	middleWindow := container.NewVBox(updateLabel, middleHBox)

	// 下方視窗; 負責顯示所有總量的資料
	curprice := canvas.NewText("", color.NRGBA{R: 0, G: 120, B: 120, A: 255})
	curprice.Alignment = fyne.TextAlignTrailing
	curprice.TextSize = 15

	topic2 := canvas.NewText("TotalAsset", color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	topic2.Alignment = fyne.TextAlignLeading
	topic2.TextSize = 20

	bottomLabel := canvas.NewText("Name: BTC       Amount: "+fmt.Sprintf("%.4f", t.Totalamount)+
		"       Cost: $"+fmt.Sprintf("%.4f", t.Totalcost)+"        Avgprice :"+fmt.Sprintf("%.4f", t.Totalavgprice), color.Black)
	bottomLabel.Alignment = fyne.TextAlignCenter
	bottomLabel.TextSize = 15

	resetButton := widget.NewButton("Reset", func() {
		t := Total{}
		err := db.QueryRow("SELECT SUM(amount) FROM cryptolist").Scan(&t.Totalamount)
		if err != nil {
			log.Println("Failed to Select SUM(amount):", err)
			return
		}
		totalAmountProift = t.Totalamount
		err = db.QueryRow("SELECT SUM(cost) FROM cryptolist").Scan(&t.Totalcost)
		if err != nil {
			log.Println("Failed to Select SUM(cost):", err)
			return
		}
		totalCost = t.Totalcost
		err = db.QueryRow("SELECT SUM(cost)/SUM(amount) FROM cryptolist").Scan(&t.Totalavgprice)
		if err != nil {
			log.Println("Failed to Select SUM(cost)/SUM(amount): ", err)
			return
		}
		bottomLabel.Text = "Name: BTC       Amount: " + fmt.Sprintf("%.4f", t.Totalamount) +
			"       Cost: $" + fmt.Sprintf("%.4f", t.Totalcost) + "        Avgprice :" + strconv.FormatFloat(t.Totalavgprice, 'f', 2, 64)
		bottomLabel.Refresh()

		log.Println("Reset Button Activated")
	})
	middleHBox.Add(resetButton)

	//總量現價
	curValue = canvas.NewText("Current Value: $0.00", color.NRGBA{R: 231, G: 171, B: 78, A: 255})
	curValue.Alignment = fyne.TextAlignLeading
	curValue.TextSize = 18

	curProfit = canvas.NewText("Total Profit: +0.00", color.NRGBA{R: 231, G: 171, B: 78, A: 255})
	curProfit.Alignment = fyne.TextAlignLeading
	curProfit.TextSize = 18

	bottomWindow := container.NewVBox(curValue, curProfit, topic2, bottomLabel)

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

// 把此次新增的交易新增到database
func insertinfo(db *sql.DB, name string, amount, cost, avgcost float64) {

	insertStmt, err := db.Prepare("INSERT INTO cryptolist(name,amount,cost,avgcost) VALUES(?,?,?,?)")
	if err != nil {
		log.Println("Failed to insert info: ", err)
		return
	}
	_, err = insertStmt.Exec(name, amount, cost, avgcost)
	if err != nil {
		log.Println("Failed to execute stmt: ", err)
		return
	}
	log.Println("Insert complete")
}

func Updatecurrentprice(price string) {
	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		log.Println("Failed to parse Price: ", err)
		return
	}
	log.Println("Current Price Updated")

	//計算總價值
	value := totalAmountProift * floatPrice
	formattedValue := strconv.FormatInt(int64(value), 10)
	formattedValueWithCommas := addCommas(formattedValue)

	profitORlost(value, totalCost)

	curPriceText.Text = "Current Price: $" + fmt.Sprintf("%.4f", floatPrice)
	curValue.Text = "Current Value: $" + formattedValueWithCommas

	curPriceText.Refresh()
	curValue.Refresh()
}

// 為總資產的數字放上標點符號
func addCommas(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return addCommas(s[:n-3]) + "," + s[n-3:]
}

func profitORlost(value, cost float64) {
	if value > cost {
		profit := value - cost
		formattedProfit := strconv.FormatInt(int64(profit), 10)
		formattedProfitWithCommas := addCommas(formattedProfit)

		percentProfitBasis := value / totalCost
		percentProfit := percentProfitBasis * 100
		formattedPercentProfit := strconv.FormatInt(int64(percentProfit), 10)
		formattedPercentProfitWithCommas := addCommas(formattedPercentProfit)
		curProfit.Text = "Total Profit: +" + formattedProfitWithCommas + "(" + formattedPercentProfitWithCommas + "%)"

		curProfit.Refresh()
	} else {
		lost := cost - value
		formattedLost := strconv.FormatInt(int64(lost), 10)
		formattedLostWithCommas := addCommas(formattedLost)

		percentLostBasis := 1 - (value / totalCost)
		percentLost := percentLostBasis * 100
		formattedPercentLost := strconv.FormatInt(int64(percentLost), 10)
		formattedPercentLostWithCommas := addCommas(formattedPercentLost)
		curProfit.Text = "Total Profit: -" + formattedLostWithCommas + "(" + formattedPercentLostWithCommas + "%)"

		curProfit.Refresh()
	}
}
