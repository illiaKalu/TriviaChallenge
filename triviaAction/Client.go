package triviaAction

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ws                  *websocket.Conn
	sendMessageChannel  chan []byte
}

func (c *Client) Write() {
	defer func() {
		c.ws.Close()
	}()

	for {
		select {

		case message, ok := <- c.sendMessageChannel:
			{

				if !ok {
					c.ws.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				c.ws.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}

func (c *Client) Read(hub *Hub) {
	defer func() {
		hub.RemoveClient <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()

		if err != nil {
			hub.RemoveClient <- c
			c.ws.Close()
			break
		}

		hub.Broadcast <- message
	}
}