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

	h *hub
)

func AnswerHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	var client *Client

	for{
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
	
		resq := make(map[string]string)
		err = json.Unmarshal(msg, &resq)
	
		if resq["type"] == "online" {
			client = h.NewClient(resq["device_id"],ws)
	
			ws.SetCloseHandler(func(code int, text string) error {
				h.Close(resq["device_id"])
				fmt.Println(resq["device_id"], " is offline,")
				return nil
			})
			
			break
		}	
	}

	go client.readPump()
	go client.writePump()
}

func OfferHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(ws)

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,websocket.CloseNoStatusReceived) {
				fmt.Printf("error: %v", err)
			}
			return
		}
	
		resq := make(map[string]string)
		err = json.Unmarshal(msg, &resq)
	
		if resq["type"] == "offer" {
			if h.ExistClient(resq["device_id"]){

				if h.client[resq["device_id"]].using {
					resp := "{\"status\":\"error\",\"msg\":\""+resq["device_id"]+" is using\"}"
					err = client.conn.WriteMessage(websocket.TextMessage, []byte(resp))
					if err != nil {
						fmt.Println(err)
						return
					}
					continue
				}
				
				h.Connection(resq["device_id"],client)

				ws.SetCloseHandler(func(code int, text string) error {
					h.Close(resq["device_id"])
					return nil
				})

				client.send_ch <- string(msg)
				data := <- client.receive_ch

				err = ws.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					fmt.Println(err)
					return
				}
				break
			}
		}
		resp := "{\"status\":\"error\",\"msg\":\""+resq["device_id"]+" is not exist\"}"
		err = client.conn.WriteMessage(websocket.TextMessage, []byte(resp))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	
	go client.readPump()
	go client.writePump()
				
}

func GetDevices(c *gin.Context) {
	devices := h.GetClients()
	c.JSON(http.StatusOK, gin.H{"data": devices})
}

func InitWSHub(){
	h = NewHub()
}