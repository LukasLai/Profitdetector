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
				log.Fatal("Unable to read message：", err)
			}

			// 解析 JSON 訊息
			var data map[string]interface{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Fatal("Unable to unmarshal json:", err)
			}

			// 提取 "b" 欄位的值
			buyPrice, ok := data["b"].(string)
			if !ok {
				log.Fatal("Unable to get buyPrice")
			}
			//暫停0.5秒
			time.Sleep(time.Millisecond * 500) //因使用並行運算，主程式若來不及創建完GUI，gui.Updatecurrentprice(buyPrice)將會出錯
			gui.Updatecurrentprice(buyPrice)
			time.Sleep(time.Second * 5)
		}
	}()
	//SQL test
	sqldata.Getsqldb()
}
