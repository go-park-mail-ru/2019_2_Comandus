package chgrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/chat_app/chat/delivery/grpc/chat_grpc"
	model2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

type Server struct {
	Ucase chat.Usecase
}

func NewChatServerGrpc(gserver *grpc.Server, ucase chat.Usecase) {
	server := &Server{
		Ucase: ucase,
	}
	chat_grpc.RegisterChatHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *Server) TransformChatRPC(chat *model2.Chat) *chat_grpc.Chat {
	if chat == nil {
		return nil
	}

	res := &chat_grpc.Chat{
		ID:                   chat.ID,
		User:                 chat.Freelancer,
		Support:              chat.Manager,
		Name:                 chat.Name,
		Proposal:             chat.ProposalId,
	}
	return res
}

func (s *Server) TransformChatData(chat *chat_grpc.Chat) *model2.Chat {
	res := &model2.Chat{
		ID:         chat.ID,
		Freelancer: chat.User,
		Manager:    chat.Support,
		Name:       chat.Name,
		ProposalId: chat.Proposal,
	}
	return res
}

func (s *Server) Create(context context.Context, chat *chat_grpc.Chat) (*chat_grpc.Chat, error) {
	log.Println("CREATE CHAT ", chat.Name)
	res, err := s.Ucase.CreateChat(s.TransformChatData(chat))
	if err != nil {
		return nil, errors.Wrap(err, "UserUcase.CreateUser")
	}
	return s.TransformChatRPC(res), nil
}

