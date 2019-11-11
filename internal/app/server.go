package apiserver

import (
	"database/sql"
	companyRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/repository"
	freelancerHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/http"
	freelancerRepository "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/repository"
	freelancerUcase "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/usecase"
	mainHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/repository"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user"
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
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

type ctxKey int8

const (
	ctxKeyUser              ctxKey = iota
	sessionName                    = "user-session"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux				*mux.Router
	db				*sql.DB
	sessionStore	sessions.Store
	config			*Config
	logger			*zap.SugaredLogger
	clientUrl		string
	token			*HashToken
	sanitizer		*bluemonday.Policy
	usecase			user.Usecase
}

func NewServer(sessionStore sessions.Store, thisLogger *zap.SugaredLogger, thisToken *HashToken, thisSanitizer *bluemonday.Policy, db *sql.DB) *server {
	s := &server{
		mux:          	mux.NewRouter(),
		sessionStore: 	sessionStore,
		logger:		  	thisLogger,
		clientUrl:    	"https://comandus.now.sh",
		token:	  	  	thisToken,
		sanitizer:	  	thisSanitizer,
		db:				db,
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) ConfigureServer() {
	userRep := userRepository.NewUserRepository(s.db)
	managerRep := managerRepository.NewManagerRepository(s.db)
	freelancerRep := freelancerRepository.NewFreelancerRepository(s.db)
	companyRep := companyRepository.NewCompanyRepository(s.db)
	jobRep := jobRepository.NewJobRepository(s.db)
	responseRep := responseRepository.NewResponseRepository(s.db)
	contractRep := contractRepository.NewContractRepository(s.db)

	userU := userUcase.NewUserUsecase(userRep, managerRep, freelancerRep)
	freelancerU := freelancerUcase.NewFreelancerUsecase(userRep, freelancerRep)
	jobU := jobUcase.NewJobUsecase(userRep, managerRep, jobRep)
	responseU := responseUcase.NewResponseUsecase(userRep, managerRep, freelancerRep, jobRep, responseRep)
	contractU := contractUcase.NewContractUsecase(userRep, managerRep, freelancerRep, jobRep, responseRep, contractRep, companyRep)

	s.usecase = userU
	mainHttp.NewMainHandler(s.mux, userU, s.sanitizer, s.logger, s.sessionStore)

	s.mux.Use(s.RequestIDMiddleware, s.CORSMiddleware, s.AccessLogMiddleware)

	// only for auth users
	private := s.mux.PathPrefix("").Subrouter()
	private.Use(s.AuthenticateUser, s.CheckTokenMiddleware)

	userHttp.NewUserHandler(private, userU, s.sanitizer, s.logger, s.sessionStore)
	freelancerHttp.NewFreelancerHandler(private, freelancerU, s.sanitizer, s.logger, s.sessionStore)
	jobHttp.NewJobHandler(private, jobU, s.sanitizer, s.logger, s.sessionStore)
	responseHttp.NewResponseHandler(private, responseU, s.sanitizer, s.logger, s.sessionStore)
	contractHttp.NewContractHandler(private, contractU, s.sanitizer, s.logger, s.sessionStore)
}

