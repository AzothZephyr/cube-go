package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/azothzephyr/cube-bot/pkg/market_data"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var apiKey string = "asdf"

type CubeBot struct {
	data []market_data.AggMessage
	ws   *websocket.Conn
}

func NewCubeBot() *CubeBot {
	// setup websocket connection
	dialer := websocket.DefaultDialer
	// dev cert is fucked up, so we hit by ip
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	conn, _, err := dialer.Dial("wss://147.28.171.25/md/tops", nil)
	// conn, _, err := dialer.Dial("wss://dev.cube.exchange/md/tops", nil)
	// conn, _, err := dialer.Dial("wss://api.cube.exchange/md/tops", nil)
	if err != nil {
		panic(err)
	}

	bot := &CubeBot{ws: conn}
	// this is blocking
	bot.run()
	return bot
}

func (bot *CubeBot) run() {
	defer bot.ws.Close()

	bot.sendCommand("auth", []string{bot.getAuth()})
	// bot.waitForAccount()
	// bot.waitForSymbol("AAPL")

	heartbeatTicker := time.NewTicker(29 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-heartbeatTicker.C:
				// do stuff
				// every 30 seconds send heart beat
				bot.sendHeartbeat()
			case <-quit:
				heartbeatTicker.Stop()
				return
			}
		}
	}()

	for {
		messageType, message, err := bot.ws.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				fmt.Printf("Websocket closed with code: %v - %v\n", closeErr.Code, closeErr.Text)
			} else {
				fmt.Println("read:", err)
			}
			return
		}

		if messageType == websocket.BinaryMessage {
			var decodedMessage market_data.AggMessage
			err := proto.Unmarshal(message, &decodedMessage)
			if err != nil {
				fmt.Println("error with message unmarshal:", err)
				continue
			}

			if decodedMessage.GetTopOfBooks() != nil {
				fmt.Println("top update")
			}
			if decodedMessage.GetHeartbeat() != nil {
				fmt.Println("heartbeat")
				hb := decodedMessage.GetHeartbeat()
				log.Println("request id:", hb.RequestId)
			}
			if decodedMessage.GetRateUpdates() != nil {
				fmt.Println("rate updates")
			}
			fmt.Println("-----")
		}
	}
}

func (bot *CubeBot) sendCommand(command string, params []string) {
	// Send command implementation

}

func (bot *CubeBot) sendHeartbeat() {
	// make random request id
	buf := make([]byte, 8)
	rand.Read(buf) // Always succeeds, no need to check error

	// create ClientMessage_Heartbeat object
	var hb market_data.ClientMessage_Heartbeat
	hb.Heartbeat = &market_data.Heartbeat{
		RequestId: binary.LittleEndian.Uint64(buf),
		Timestamp: uint64(time.Now().Unix()),
	}

	// instantiate client message object
	var cm market_data.ClientMessage
	// wrap heartbeat object in client message object
	cm.Inner = &hb

	msg, err := proto.Marshal(&cm)
	if err != nil {
		fmt.Println("Error marshalling the heartbeat message: ", err)
		return
	}
	err = bot.ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		fmt.Println("Error sending the heartbeat message: ", err)
	}
}

// func (bot *CubeBot) reconnect() {
// 	for {
// 		// Try to reconnect with a backoff
// 		time.Sleep(time.Second * 5) // For example, wait for 5 seconds before trying to reconnect
// 		err := bot.connect()        // Assume bot.connect() is a method that tries to establish a WebSocket connection
// 		if err == nil {
// 			break // Break out of the loop if connected successfully
// 		}
// 		fmt.Println("Reconnect failed:", err)
// 	}
// }

func (bot *CubeBot) waitForAccount() {
	// Implementation here
}

func (bot *CubeBot) waitForSymbol(symbol string) {
	// Implementation here
}

func (bot *CubeBot) getAuth() string {
	secretKey := []byte(apiKey)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	message := append([]byte("cube.xyz"), int64ToBytes(timestamp)...)
	signature := hmac.New(sha256.New, secretKey).Sum(message)
	return base64.StdEncoding.EncodeToString(signature)
}

func int64ToBytes(i int64) []byte {
	bytes := make([]byte, 8)
	for index := range bytes {
		bytes[index] = byte(i)
		i >>= 8
	}
	return bytes
}

func main() {
	// rand.Seed(time.Now().UnixNano())
	// instantiate cubeorderbook
	NewCubeBot()
	// while true {

	// }
}
