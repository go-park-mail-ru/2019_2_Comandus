package server_clients

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	"log"
)

type ServerClients struct {
	AuthClient			*clients.AuthClient
	CompanyClient		*clients.CompanyClient
	FreelancerClient	*clients.FreelancerClient
	JobClient			*clients.JobClient
	ManagerClient		*clients.ManagerClient
	ResponseClient		*clients.ResponseClient
	UserClient			*clients.UserClient
}

func NewClients() *ServerClients {
	sc := &ServerClients{
		AuthClient:       new(clients.AuthClient),
		CompanyClient:    new(clients.CompanyClient),
		FreelancerClient: new(clients.FreelancerClient),
		JobClient:        new(clients.JobClient),
		ManagerClient:    new(clients.ManagerClient),
		ResponseClient:   new(clients.ResponseClient),
		UserClient:       new(clients.UserClient),
	}

	if err := sc.AuthClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.FreelancerClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.CompanyClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.JobClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.ManagerClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.ResponseClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := sc.UserClient.Connect(); err != nil {
		log.Println(err)
	}

	return sc
}
