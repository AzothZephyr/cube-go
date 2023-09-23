package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
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

	heartbeatTicker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-heartbeatTicker.C:
				// do stuff
				// every 30 seconds send heart beat
				fmt.Println("sending heartbeat")
				var hb market_data.Heartbeat
				hb.RequestId = 1
				hb.Timestamp = uint64(time.Now().Unix())

				bot.sendCommand(hb.String(), []string{bot.getAuth()})
			case <-quit:
				heartbeatTicker.Stop()
				return
			}
		}
	}()

	for {
		messageType, message, err := bot.ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		if messageType == websocket.BinaryMessage {
			var decodedMessage market_data.AggMessage
			err := proto.Unmarshal(message, &decodedMessage)
			if err != nil {
				fmt.Println("error with message unmarshal:", err)
				continue
			}
			// fmt.Println(decodedMessage.GetTopOfBooks())

			if decodedMessage.GetTopOfBooks() != nil {
				fmt.Println("top update")
				// decodedMessage.Inner.
			}
			if decodedMessage.GetHeartbeat() != nil {
				fmt.Println("heartbeat")
			}
			if decodedMessage.GetRateUpdates() != nil {
				fmt.Println("rate updates")
			}
			// msgType := reflect.TypeOf(&decodedMessage).Kind()
			// fmt.Println(decodedMessage.ProtoReflect().Type())
			fmt.Println("-----")
			// switch string(msgType) {
			// default:
			// 	fmt.Printf("unexpected type %T", v)
			// case "market_data.AggMessage_Heartbeat":
			// 	e.code = Code(C.curl_wrapper_easy_setopt_long(e.curl, C.CURLoption(option), C.long(v)))
			// case string:
			// 	e.code = Code(C.curl_wrapper_easy_setopt_str(e.curl, C.CURLoption(option), C.CString(v)))
			// }
		}
	}
}

func (bot *CubeBot) sendCommand(command string, params []string) {
	// Send command implementation
}

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
	rand.Seed(time.Now().UnixNano())
	// instantiate cubeorderbook
	NewCubeBot()
	// while true {

	// }
}
