package apiserver

import (
	"database/sql"
	coRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company/repository"
	http3 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/delivery/http"
	fRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/repository"
	usecase2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/freelancer/usecase"
	generalHttp "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/delivery/http"
	mRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/repository"
	http5 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/delivery/http"
	cRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/repository"
	usecase5 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract/usecase"
	http4 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/delivery/http"
	jRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/repository"
	usecase3 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-job/usecase"
	http6 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/delivery/http"
	rRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/repository"
	usecase4 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-response/usecase"
	http2 "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/delivery/http"
	uRep "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/repository"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user/usecase"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store"
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
	hireManagerIdCookieName        = "hire-manager-id"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	mux				*mux.Router
	store			store.Store
	db				*sql.DB
	sessionStore	sessions.Store
	config			*Config
	logger			*zap.SugaredLogger
	clientUrl		string
	token			*HashToken
	sanitizer		*bluemonday.Policy
}

func NewServer(sessionStore sessions.Store, store store.Store, thisLogger *zap.SugaredLogger, thisToken *HashToken, thisSanitizer *bluemonday.Policy) *server {
	s := &server{
		mux:          mux.NewRouter(),
		sessionStore: sessionStore,
		logger:		  thisLogger,
		clientUrl:    "https://comandus.now.sh",
		store:        store,
		token:	  	  thisToken,
		sanitizer:	  thisSanitizer,
	}
	s.ConfigureServer()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) ConfigureServer() {
	userRep := uRep.NewUserRepository(s.db)
	managerRep := mRep.NewManagerRepository(s.db)
	freelancerRep := fRep.NewFreelancerRepository(s.db)
	companyRep := coRep.NewCompanyRepository(s.db)
	jobRep := jRep.NewJobRepository(s.db)
	responseRep := rRep.NewResponseRepository(s.db)
	contractRep := cRep.NewContractRepository(s.db)

	userUcase := usecase.NewUserUsecase(userRep, managerRep, freelancerRep)
	freelancerUcase := usecase2.NewFreelancerUsecase(userRep, freelancerRep)
	jobUcase := usecase3.NewJobUsecase(userRep, managerRep, jobRep)
	responseUcase := usecase4.NewResponseUsecase(userRep, managerRep, freelancerRep, jobRep, responseRep)
	contractUcase := usecase5.NewContractUsecase(userRep, managerRep, freelancerRep, jobRep, responseRep, contractRep, companyRep)

	generalHttp.NewMainHandler(s.mux, userUcase, s.sanitizer, s.logger, s.sessionStore)

	s.mux.Use(s.RequestIDMiddleware, s.CORSMiddleware, s.AccessLogMiddleware)
	private := s.mux.PathPrefix("").Subrouter()
	private.Use(s.AuthenticateUser, s.CheckTokenMiddleware)

	http2.NewUserHandler(private, userUcase, s.sanitizer, s.logger, s.sessionStore)
	http3.NewFreelancerHandler(private, freelancerUcase, s.sanitizer, s.logger, s.sessionStore)
	http4.NewJobHandler(private, jobUcase, s.sanitizer, s.logger, s.sessionStore)
	http5.NewContractHandler(private, contractUcase, s.sanitizer, s.logger, s.sessionStore)
	http6.NewResponseHandler(private, responseUcase, s.sanitizer, s.logger, s.sessionStore)
}

