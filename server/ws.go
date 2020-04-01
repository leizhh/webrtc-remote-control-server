package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func AnswerHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("%s\n", msg)

	resp := make(map[string]string)
	err = json.Unmarshal(msg, &resp)
	var device chan string

	if resp["type"] == "online" {
		device = DeviceChan.NewChan(resp["device_id"])
		ws.SetCloseHandler(func(code int, text string) error {
			DeviceChan.Delete(resp["device_id"])
			device <- "device has been offline"
			fmt.Println(resp["device_id"], " is offline,")
			return nil
		})
	}

	for {
		data := <- device
		err = ws.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		device <- string(msg)
	}
}

func OfferHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	var device chan string

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			  return
		}
		//fmt.Printf("%s\n", msg)
	
		resp := make(map[string]string)
		err = json.Unmarshal(msg, &resp)
	
		if resp["type"] == "offer" {
			if DeviceChan.Exist(resp["device_id"]){
				device = DeviceChan.GetChan(resp["device_id"])
				device <- resp["sdp"]	
				data := <- device
				err = ws.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					fmt.Println(err)
					return
				}
				break
			}
		}
		err = ws.WriteMessage(websocket.TextMessage, []byte("device_id is not exist"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		device <- string(msg)
		data := <- device

		err = ws.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func GetDevices(c *gin.Context) {
	devices := DeviceChan.GetKeys()
	c.JSON(http.StatusOK, gin.H{"data": devices})
}
