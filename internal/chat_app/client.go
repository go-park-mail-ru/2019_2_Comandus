package chat_app

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model/chat"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	id         	int
	userId     	int64
	clientId   	int64
	senderId 	int64
	receiverId	int64
	chatId		int64
	ws         	*websocket.Conn
	server     	*Server
	ch         	chan *model.Message
	doneCh     	chan bool
}

// Create new app client.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil || server == nil {
		panic("ws or server cannot be nil")
	}

	maxId++
	ch := make(chan *model.Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{
		id:       maxId,
		userId:   0,
		clientId: 0,
		ws:       ws,
		server:   server,
		ch:       ch,
		doneCh:   doneCh,
	}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *model.Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		c.server.Err(fmt.Errorf("client %d is disconnected.", c.id))
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var input model.Packet

			err := websocket.JSON.Receive(c.ws, &input)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.errCh <- err
			} else {
				if input.Transaction == "init" {

					log.Println(input)
					c.userId = input.Chat.UserID
					c.clientId = input.Chat.SupportID

					currChat, err := c.server.ChatUcase.CreateChat(&input.Chat)
					if err != nil && err.Error() == chat.CHAT_CONFLICT_ERR {
						c.server.errCh <- err
					}

					c.chatId = currChat.ID

					log.Println("userID ", c.userId)
					log.Println("chatID ", currChat.ID)
					messages, err := c.server.MesUcase.ListByUser(c.userId, currChat.ID)
					if err != nil {
						c.server.errCh <- err
					}
					log.Println(messages)

					for _, msg := range messages {
						c.server.sendAll(msg)
					}
				} else if input.Transaction == "mes" {
					msg := input.Message
					if err := c.server.MesUcase.Create(&msg); err != nil {
						c.server.errCh <- err
					}
				} else {
					messages, err := c.server.MesUcase.ListByUser(c.userId, c.chatId)
					if err != nil {
						c.server.errCh <- err
					}
					log.Println(messages)
				}
			}
		}
	}
}