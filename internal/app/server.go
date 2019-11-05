package apiserver

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store/sqlstore"
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
	mux          *mux.Router
	store        *sqlstore.Store
	sessionStore sessions.Store
	config       *Config
	logger    	 *zap.SugaredLogger
	clientUrl    string
	token 	 	 *HashToken
	sanitizer	 *bluemonday.Policy
}

func newServer(sessionStore sessions.Store, store *sqlstore.Store, thisLogger *zap.SugaredLogger,
	thisToken *HashToken, thisSanitizer *bluemonday.Policy) *server {
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
	s.mux.HandleFunc("/", s.HandleMain)
	s.mux.HandleFunc("/signup", s.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	s.mux.HandleFunc("/login", s.HandleSessionCreate).Methods(http.MethodPost, http.MethodOptions)
	s.mux.Use(s.RequestIDMiddleware, s.CORSMiddleware, s.AccessLogMiddleware)

	// only for authenticated users
	private := s.mux.PathPrefix("").Subrouter()
	private.Use(s.authenticateUser, s.CheckTokenMiddleware)
	private.HandleFunc("/token" , s.HandleGetToken).Methods(http.MethodGet)

	private.HandleFunc("/setusertype", s.HandleSetUserType).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/account", s.HandleShowProfile).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account", s.HandleEditProfile).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/upload-avatar", s.HandleUploadAvatar).Methods(http.MethodPost, http.MethodOptions)

	private.HandleFunc("/account/download-avatar", s.HandleDownloadAvatar).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/account/avatar/{id:[0-9]}", s.HandleGetAvatar).Methods(http.MethodGet, http.MethodOptions)

	private.HandleFunc("/account/settings/password", s.HandleEditPassword).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/account/settings/notifications", s.HandleEditNotifications).Methods(http.MethodPut, http.MethodOptions)
	//private.HandleFunc("/account/settings/auth-history", s.HandleGetAuthHistory).Methods(http.MethodGet, http.MethodOptions)

	private.HandleFunc("/roles", s.HandleRoles).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodDelete, http.MethodOptions)
	private.HandleFunc("/jobs", s.HandleCreateJob).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/jobs", s.HandleGetAllJobs).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/jobs/{id:[0-9]+}", s.HandleGetJob).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/jobs/{id:[0-9]+}", s.HandleUpdateJob).Methods(http.MethodPut, http.MethodOptions)

	private.HandleFunc("/jobs/{id:[0-9]+}/response}", s.HandleResponseJob).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/responses", s.HandleGetResponses).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/responses/{id:[0-9]+}/accept", s.HandleResponseAccept).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/responses/{id:[0-9]+}/deny", s.HandleResponseDeny).Methods(http.MethodPut, http.MethodOptions)

	private.HandleFunc("/responses/{id:[0-9]+}/contract", s.HandleCreateContract).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/contract/{id:[0-9]+/done}", s.HandleTickContractAsDone).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/contract/{id:[0-9]+}", s.HandleReviewContract).Methods(http.MethodPut, http.MethodOptions)

	private.HandleFunc("/freelancer", s.HandleEditFreelancer).Methods(http.MethodPut, http.MethodOptions)
	private.HandleFunc("/freelancer/{freelancerId}", s.HandleGetFreelancer).Methods(http.MethodGet, http.MethodOptions)
}

func (s *server) HandleMain(w http.ResponseWriter, r *http.Request) {
	s.respond(w, r, http.StatusOK, "hello from server")
}

// error handler
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	ctx := r.Context()
	reqID := ctx.Value("rIDKey").(string)
	s.logger.Infof("Request ID: %s | error : %s", reqID , err.Error())
	s.respond(w, r, code, map[string]string{"error": errors.Cause(err).Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
