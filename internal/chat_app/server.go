package chat_app

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/message"
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	pattern		string
	clients		map[int]*Client
	addCh		chan *Client
	delCh		chan *Client
	sendAllCh	chan *model2.Message
	doneCh		chan bool
	errCh		chan error
	initCh		chan bool
	MesUcase	message.Usecase
	ChatUcase	chat.Usecase
}

func NewServer(pattern string, mUcase message.Usecase, chUcase chat.Usecase) *Server {
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *model2.Message)
	doneCh := make(chan bool)
	errCh := make(chan error)
	initCh := make(chan bool)

	return &Server{
		pattern,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
		initCh,
		mUcase,
		chUcase,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *model2.Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	messages, _ := s.MesUcase.List(c.chatId)
	for _, msg := range messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *model2.Message) {
	log.Println("send All:", msg)
	for _, c := range s.clients {
		if c.chatId == msg.ChatID {
			c.Write(msg)
		}
	}
}

func (s *Server) Listen() {
	log.Println("Listening server...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))

	log.Println("Created handler")
	for {
		select {

		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			//s.sendPastMessages(c)

		case c := <-s.delCh:
			delete(s.clients, c.id)

		case msg := <-s.sendAllCh:
			if err := s.MesUcase.Create(msg); err != nil {
				log.Println(err)
			}
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}