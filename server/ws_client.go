package server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct{
	conn *websocket.Conn
	receive_ch chan string
	send_ch chan string
	using bool
}

func NewClient(ws *websocket.Conn)* Client{
	client := &Client{
		conn:ws,
		receive_ch:nil,
		send_ch:nil,
		using:false,
	}

	return client
}

func (c *Client)readPump(){
	defer c.conn.Close()

	for{
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,websocket.CloseNoStatusReceived) {
				fmt.Printf("error: %v", err)
			}
			return
		}
		c.send_ch <- string(msg)
	}
}

func (c *Client)writePump(){
	defer c.conn.Close()

	for data := range c.receive_ch{
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				fmt.Println(err)
				return
			}
	}
	err := c.conn.WriteMessage(websocket.TextMessage, []byte("connection has been closed"))
	if err != nil {
		fmt.Println(err)
		return
	}
}