package apiserver

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients"
	cogrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/grpc"
	companyHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/http"
	companyRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/repository"
	companyUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/usecase"
	fgrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/grpc"
	freelancerHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/http"
	freelancerRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/repository"
	freelancerUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/usecase"
	authgrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/grpc"
	mainHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	generalUsecase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/usecase"
	mgrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/repository"
	managerUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	contractHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/delivery/http"
	contractRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/repository"
	contractUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/usecase"
	jgrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/grpc"
	jobHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/http"
	jobRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/repository"
	jobUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/usecase"
	regrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/grpc"
	responseHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/http"
	responseRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/repository"
	responseUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/usecase"
	ugrpc "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/grpc"
	userHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/repository"
	userUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type Server struct {
	Mux          *mux.Router
	SessionStore sessions.Store
	Config       *Config
	Logger       *zap.SugaredLogger
	Token        *token.HashToken
	Sanitizer    *bluemonday.Policy
}

func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	hashToken, err := token.NewHMACHashToken(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Mux:          mux.NewRouter(),
		SessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
		Logger:       logger,
		Token:        hashToken,
		Sanitizer:    bluemonday.UGCPolicy(),
		Config:       config,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

// TODO: separate function
func (s *Server) ConfigureServer(db *sql.DB) {

	userRep := userRepository.NewUserRepository(db)
	managerRep := managerRepository.NewManagerRepository(db)
	freelancerRep := freelancerRepository.NewFreelancerRepository(db)
	companyRep := companyRepository.NewCompanyRepository(db)
	jobRep := jobRepository.NewJobRepository(db)
	responseRep := responseRepository.NewResponseRepository(db)
	contractRep := contractRepository.NewContractRepository(db)

	userClient := new(clients.UserClient)
	freelancerClient := new(clients.FreelancerClient)
	managerClient := new(clients.ManagerClient)
	companyClient := new(clients.CompanyClient)
	jobClient := new(clients.JobClient)
	responseClient := new(clients.ResponseClient)
	authClient := new(clients.AuthClient)

	if err := authClient.Connect(); err != nil {
		log.Println(err)
	}
	if err := freelancerClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := companyClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := jobClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := managerClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := responseClient.Connect(); err != nil {
		log.Println(err)
	}

	if err := userClient.Connect(); err != nil {
		log.Println(err)
	}


	userU := userUcase.NewUserUsecase(userRep, freelancerClient, managerClient, companyClient)
	companyU := companyUcase.NewCompanyUsecase(companyRep, managerClient)
	managerU := managerUcase.NewManagerUsecase(managerRep)
	freelancerU := freelancerUcase.NewFreelancerUsecase(freelancerRep)
	jobU := jobUcase.NewJobUsecase(jobRep, managerClient)
	responseU := responseUcase.NewResponseUsecase(responseRep, freelancerClient, managerClient, jobClient)
	contractU := contractUcase.NewContractUsecase(contractRep, freelancerClient, managerClient, companyClient, jobClient, responseClient)
	generalU := generalUsecase.NewGeneralUsecase(authClient, freelancerClient, managerClient, companyClient)

	s.Mux.Handle("/metrics", promhttp.Handler())
	private := s.Mux.PathPrefix("").Subrouter()

	mainHttp.NewMainHandler(s.Mux, private, s.Sanitizer, s.Logger, s.SessionStore, s.Token, generalU)

	mid := middleware.NewMiddleware(s.SessionStore, s.Logger, s.Token, s.Config.ClientUrl, userClient)
	s.Mux.Use(mid.RequestIDMiddleware, mid.CORSMiddleware, mid.AccessLogMiddleware)

	// only for auth users
	private.Use(mid.AuthenticateUser, mid.CheckTokenMiddleware)

	userHttp.NewUserHandler(private, userU, s.Sanitizer, s.Logger, s.SessionStore)
	freelancerHttp.NewFreelancerHandler(s.Mux, private, freelancerU, s.Sanitizer, s.Logger, s.SessionStore)
	jobHttp.NewJobHandler(s.Mux, private, jobU, s.Sanitizer, s.Logger, s.SessionStore)
	companyHttp.NewCompanyHandler(private, companyU, s.Sanitizer, s.Logger, s.SessionStore)
	responseHttp.NewResponseHandler(private, responseU, s.Sanitizer, s.Logger, s.SessionStore)
	contractHttp.NewContractHandler(private, contractU, s.Sanitizer, s.Logger, s.SessionStore)

	go func() {
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		authgrpc.NewAuthServerGrpc(server, userU)

		fmt.Println("starting server at :8081")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8087")
		if err != nil {
			log.Fatalln("cant listen port", err)
		}
		server := grpc.NewServer()
		ugrpc.NewUserServerGrpc(server, userU)

		fmt.Println("starting server at :8087")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		cogrpc.NewCompanyServerGrpc(server, companyU)

		fmt.Println("starting server at :8082")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8083")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		fgrpc.NewFreelancerServerGrpc(server, freelancerU)

		fmt.Println("starting server at :8083")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8084")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		mgrpc.NewManagerServerGrpc(server, managerU)

		fmt.Println("starting server at :8084")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8085")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		jgrpc.NewJobServerGrpc(server, jobU)

		fmt.Println("starting server at :8085")
		server.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", ":8086")
		if err != nil {
			log.Fatalln("cant listet port", err)
		}
		server := grpc.NewServer()
		regrpc.NewResponseServerGrpc(server, responseU)

		fmt.Println("starting server at :8086")
		server.Serve(lis)
	}()


}

