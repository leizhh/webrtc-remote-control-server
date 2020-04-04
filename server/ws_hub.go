package server

import (
	"github.com/gorilla/websocket"
)

type hub struct{
	client map[string]*Client
}

func (h *hub)NewClient(key string,ws *websocket.Conn)* Client{
	client := &Client{
		conn:ws,
		receive_ch:make(chan string,1),
		send_ch:make(chan string,1),
		using:false,
	}
	h.client[key]=client

	return client
}

func (h *hub)GetClient(key string)* Client{
	return h.client[key]
}

func (h *hub) ExistClient(key string)bool{
	if _, ok := h.client[key]; ok {
		return true
	}
	return false
}


type resp_clients struct {
	Device_id string `json:"device_id"`
	Using bool `json:"using"`
}

func (h *hub) GetClients()[]resp_clients{
	var res []resp_clients

	for k , v  := range h.client {
		tmp := resp_clients{
			k,
			v.using,
		}
		res = append(res,tmp)
	}
	return res
}

func (h *hub) Close(key string){
	close(h.client[key].send_ch)
	close(h.client[key].receive_ch)
	delete(h.client, key)
}


func (h *hub) Connection(key string,client *Client){
	client.receive_ch = h.client[key].send_ch
	client.send_ch = h.client[key].receive_ch
	h.client[key].using = true
	client.using = true
}

func NewHub()* hub{
	 return &hub{
		make(map[string]*Client),
	}
}