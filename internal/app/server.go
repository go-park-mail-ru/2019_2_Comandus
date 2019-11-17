package apiserver

import (
	"database/sql"
	companyHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/http"
	companyRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/repository"
	companyUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/usecase"
	freelancerHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/http"
	freelancerRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/repository"
	freelancerUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/usecase"
	mainHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/repository"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	contractHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/delivery/http"
	contractRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/repository"
	contractUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/usecase"
	jobHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/http"
	jobRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/repository"
	jobUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/usecase"
	responseHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/http"
	responseRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/repository"
	responseUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/usecase"
	userHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/repository"
	userUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
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

	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep, companyRep)
	companyU := companyUcase.NewCompanyUsecase(companyRep, managerRep)
	freelancerU := freelancerUcase.NewFreelancerUsecase(freelancerRep)
	jobU := jobUcase.NewJobUsecase(managerRep, jobRep)
	responseU := responseUcase.NewResponseUsecase(managerRep, freelancerRep, jobRep, responseRep)
	contractU := contractUcase.NewContractUsecase(managerRep, freelancerRep, jobRep, responseRep, contractRep)

	private := s.Mux.PathPrefix("").Subrouter()

	mainHttp.NewMainHandler(s.Mux, private, userU, s.Sanitizer, s.Logger, s.SessionStore, s.Token)

	mid := middleware.NewMiddleware(s.SessionStore, s.Logger, s.Token, userU, s.Config.ClientUrl)
	s.Mux.Use(mid.RequestIDMiddleware, mid.CORSMiddleware, mid.AccessLogMiddleware)

	// only for auth users
	private.Use(mid.AuthenticateUser, mid.CheckTokenMiddleware)

	userHttp.NewUserHandler(private, userU, s.Sanitizer, s.Logger, s.SessionStore)
	freelancerHttp.NewFreelancerHandler(private, freelancerU, s.Sanitizer, s.Logger, s.SessionStore)
	jobHttp.NewJobHandler(private, jobU, s.Sanitizer, s.Logger, s.SessionStore)
	companyHttp.NewCompanyHandler(private, companyU, s.Sanitizer, s.Logger, s.SessionStore)
	responseHttp.NewResponseHandler(private, responseU, s.Sanitizer, s.Logger, s.SessionStore)
	contractHttp.NewContractHandler(private, contractU, s.Sanitizer, s.Logger, s.SessionStore)
}
