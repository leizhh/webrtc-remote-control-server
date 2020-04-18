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

type Session struct {
	Type     string `json:"type"`
	DeviceId string `json:"device_id"`
	Msg      string `json:"msg"`
	Data     string `json:"data"`
}

func AnswerHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	var client *Client

	for {
		req := &Session{}
		if err := ws.ReadJSON(req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				fmt.Printf("error: %v", err)
			}
			return
		}

		if req.Type == "online" {
			if h.ExistClient(req.DeviceId) {
				resp := Session{}
				resp.Type = "error"
				resp.Msg = req.DeviceId + " is exist"
				err = ws.WriteJSON(resp)
				continue
			}

			client = h.NewClient(req.DeviceId, ws)

			client.conn.SetCloseHandler(func(code int, text string) error {
				h.Close(req.DeviceId)
				fmt.Println(req.DeviceId, " is offline,")
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
		req := &Session{}
		if err := ws.ReadJSON(req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				fmt.Printf("error: %v", err)
			}
			return
		}

		if req.Type == "offer" {
			if h.ExistClient(req.DeviceId) {
				if h.client[req.DeviceId].using {
					resp := Session{}
					resp.Type = "error"
					resp.Msg = req.DeviceId + " is using"

					if err := ws.WriteJSON(resp); err != nil {
						fmt.Println(err)
						return
					}
					continue
				}

				h.Connection(req.DeviceId, client)

				client.conn.SetCloseHandler(func(code int, text string) error {
					h.Close(req.DeviceId)
					return nil
				})

				msg, _ := json.Marshal(req)
				client.send_ch <- string(msg)
				data := <-client.receive_ch

				if err := client.conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
					fmt.Println(err)
					return
				}
				break
			}
		}

		resp := Session{}
		resp.Type = "error"
		resp.Msg = req.DeviceId + " is not exist"
		if err := client.conn.WriteJSON(resp); err != nil {
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

func InitDeviceHub() {
	h = NewHub()
}
