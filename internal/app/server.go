package apiserver

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	companyHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/delivery/http"
	companyRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/repository"
	companyUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer"
	freelancerHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/http"
	freelancerRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/repository"
	freelancerUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/usecase"
	mainHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/repository"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	contractHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/delivery/http"
	contractRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/repository"
	contractUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/usecase"
	user_job "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job"
	jobHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/http"
	jobRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/repository"
	jobUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/usecase"
	userresponse "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response"
	responseHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/http"
	responseRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/repository"
	responseUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/usecase"
	userHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/repository"
	userUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ctxKey int8

const (
	CtxKeyUser  ctxKey = iota
	sessionName        = "user-session"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type Server struct {
	mux          *mux.Router
	sessionStore sessions.Store
	config       *Config
	logger       *zap.SugaredLogger
	clientUrl    string
	token        *HashToken
	sanitizer    *bluemonday.Policy
	usecase      user.Usecase
}

type ServerRepos struct {
	userRep       user.Repository
	managerRep    manager.Repository
	freelancerRep freelancer.Repository
	companyRep    company.Repository
	jobRep        user_job.Repository
	responseRep   userresponse.Repository
	contractRep   user_contract.Repository
}

func (sr *ServerRepos) NewRepos(db *sql.DB) {
	sr.userRep = userRepository.NewUserRepository(db)
	sr.managerRep = managerRepository.NewManagerRepository(db)
	sr.freelancerRep = freelancerRepository.NewFreelancerRepository(db)
	sr.companyRep = companyRepository.NewCompanyRepository(db)
	sr.jobRep = jobRepository.NewJobRepository(db)
	sr.responseRep = responseRepository.NewResponseRepository(db)
	sr.contractRep = contractRepository.NewContractRepository(db)
}

func NewServer(m *mux.Router, sessionStore sessions.Store, thisLogger *zap.SugaredLogger, thisToken *HashToken, thisSanitizer *bluemonday.Policy) *Server {
	s := &Server{
		mux:          m,
		sessionStore: sessionStore,
		logger:       thisLogger,
		clientUrl:    "https://comandus.now.sh",
		token:        thisToken,
		sanitizer:    thisSanitizer,
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
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

	s.usecase = userU
	mainHttp.NewMainHandler(s.mux, userU, s.sanitizer, s.logger, s.sessionStore)

	s.mux.Use(s.RequestIDMiddleware, s.CORSMiddleware, s.AccessLogMiddleware)

	// only for auth users
	private := s.mux.PathPrefix("").Subrouter()
	private.Use(s.AuthenticateUser) //, s.CheckTokenMiddleware)

	userHttp.NewUserHandler(private, userU, s.sanitizer, s.logger, s.sessionStore)
	freelancerHttp.NewFreelancerHandler(private, freelancerU, s.sanitizer, s.logger, s.sessionStore)
	companyHttp.NewCompanyHandler(private, companyU, s.sanitizer, s.logger, s.sessionStore)
	jobHttp.NewJobHandler(private, jobU, s.sanitizer, s.logger, s.sessionStore)
	responseHttp.NewResponseHandler(private, responseU, s.sanitizer, s.logger, s.sessionStore)
	contractHttp.NewContractHandler(private, contractU, s.sanitizer, s.logger, s.sessionStore)
}
