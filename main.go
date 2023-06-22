package main

import (
	"encoding/json"
	"log"
	sqldata "profitdetector/SQLdata"
	gui "profitdetector/fynegui"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	endpoint := "wss://stream.binance.com:9443/ws/btcusdt@ticker"

	// 建立 websocket 連線
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatal("無法建立 websocket 連線：", err)
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("無法讀取訊息：", err)
			}

			// 解析 JSON 訊息
			var data map[string]interface{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Fatal("無法解析 JSON 訊息：", err)
			}

			// 提取 "b" 欄位的值
			buyPrice, ok := data["b"].(string)
			if !ok {
				log.Fatal("無法取得 buyPrice")
			}

			gui.Updatecurrentprice(buyPrice)
			time.Sleep(time.Second * 5)
		}
	}()
	//SQL test
	sqldata.Getsqldb()
}
