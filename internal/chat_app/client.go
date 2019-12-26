package chat_app

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100
var maxId int = 0

type Client struct {
	id         		int
	userId			int64
	//isFreelancer	bool
	//chatId			int64
	ws         		*websocket.Conn
	server     		*Server
	ch         		chan *model.Message
	doneCh     		chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil || server == nil {
		panic("ws or server cannot be nil")
	}

	maxId++
	ch := make(chan *model.Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{
		id:       maxId,
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
		c.server.Err(fmt.Errorf("client %d is disconnected", c.id))
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		case msg := <-c.ch:
			log.Println("Send:", msg)
			if err := websocket.JSON.Send(c.ws, msg); err != nil {
				log.Println(err)
			}

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return

		default:
			var input model.Packet
			err := websocket.JSON.Receive(c.ws, &input)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.errCh <- err
			} else {
				log.Println(input)
				if input.Transaction == "init" {
					c.initChat(input)
				} else if input.Transaction == "mes" {
					c.sendMes(input)
				} else {
					log.Println("wrong request")
				}
			}
		}
	}
}

func (c *Client) initChat(input model.Packet) {
	currChat, err := c.server.ChatUcase.FindByProposal(input.Chat.ProposalId)
	if err != nil {
		c.server.errCh <- errors.New("no access")
		return
	}

	c.userId = input.Chat.UserId

	if input.Client {
		err := c.server.MesUcase.UpdateStatus(currChat.ID, currChat.Manager)
		if err != nil {
			c.server.errCh <- err
			return
		}
	} else {
		err := c.server.MesUcase.UpdateStatus(currChat.ID, currChat.Freelancer)
		if err != nil {
			c.server.errCh <- err
			return
		}
	}

	//c.chatId = currChat.ID
	messages, err := c.server.MesUcase.List(currChat.ID)
	if err != nil {
		c.server.errCh <- err
	}

	for _, msg := range messages {
		c.server.sendAll(msg)
	}
}

func (c *Client) sendMes(input model.Packet) {
	msg := input.Message
	msg.SenderID = c.userId

	chat, err := c.server.ChatUcase.FindByProposal(msg.ProposalId)
	if err != nil {
		c.server.errCh <- err
		return
	}

	msg.ChatID = chat.ID

	if c.userId != chat.Manager {
		msg.ReceiverID = chat.Manager
	} else {
		msg.ReceiverID = chat.Freelancer
	}

	for _, client := range c.server.clients {
		if client.userId == msg.SenderID || client.userId == msg.ReceiverID {
			msg.IsRead = true
			break
		}
	}

	if err := c.server.MesUcase.Create(&msg); err != nil {
		c.server.errCh <- err
		return
	}
	c.server.sendAll(&msg)
}