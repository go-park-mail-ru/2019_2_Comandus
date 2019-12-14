package mainHttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/token"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type MainHandler struct {
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
	token 			*token.HashToken
	ucase			general.Usecase
}

func NewMainHandler(m *mux.Router,private *mux.Router, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger,
	sessionStore sessions.Store, thisToken *token.HashToken, ucase general.Usecase) {
		handler := &MainHandler{
		sanitizer:		sanitizer,
		logger:			logger,
		sessionStore:	sessionStore,
		token:			thisToken,
		ucase:			ucase,
	}

	m.HandleFunc("/", handler.HandleMain)
	m.HandleFunc("/signup", handler.HandleCreateUser).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/login", handler.HandleSessionCreate).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/suggest", handler.HandleGetSuggest).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/token", handler.HandleGetToken).Methods(http.MethodGet)
	private.HandleFunc("/logout", handler.HandleLogout).Methods(http.MethodDelete, http.MethodOptions)
}

func (h *MainHandler) HandleMain(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.Labels{"path":"/", "method":"no"}))
	defer timer.ObserveDuration()

	respond.Respond(w, r, http.StatusOK, "hello from server")
}

func (h *MainHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/signup", "method":r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser<-Body.Close")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	newUser := new(model.User)
	if err := newUser.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-user.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	createdUser, err := h.ucase.CreateUser(newUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-ucase.CreateUser()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	session, err := h.sessionStore.Get(r, respond.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionGet()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = createdUser.ID
	if err := h.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionSave()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	log.Println(createdUser)
	newUser.Sanitize(h.sanitizer)
	respond.Respond(w, r, http.StatusCreated, createdUser)
}

func (h * MainHandler) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/login", "method":r.Method}))
	defer timer.ObserveDuration()

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSessionCreate<-rBodyClose()")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-ioutil.ReadAll()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	currUser := new(model.User)
	if err := currUser.UnmarshalJSON(body); err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-currCompany.UnmarshalJSON()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.ucase.VerifyUser(currUser)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-ucase.VerifyUser()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	session, err := h.sessionStore.Get(r, respond.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionGet()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = id
	if err := h.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionSave()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *MainHandler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"token", "method":r.Method}))
	defer timer.ObserveDuration()

	sess, err := h.sessionStore.Get(r, respond.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetToken<-sessionStore.Get()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	currToken, err := h.token.Create(sess, time.Now().Add(24*time.Hour).Unix())
	respond.Respond(w, r, http.StatusOK, map[string]string{"csrf-token": currToken})
}

func (h * MainHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
		Labels{"path":"/logout", "method":r.Method}))
	defer timer.ObserveDuration()

	session, err := h.sessionStore.Get(r, respond.SessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionGet()")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionSave()")
		respond.Error(w, r, http.StatusExpectationFailed, err)
		return
	}
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h *MainHandler) HandleGetSuggest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timer := prometheus.NewTimer(monitoring.RequestDuration.With(prometheus.
	Labels{"path":"suggest", "method":r.Method}))
	defer timer.ObserveDuration()

	var updateService bool
	lastReq, ok := r.Context().Value("prev_req").(string)
	if !ok || !strings.Contains(lastReq, "suggest"){
		updateService = true
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		respond.Respond(w, r, http.StatusOK, nil)
	}

	dict := r.URL.Query().Get("dict")
	if dict == "" {
		respond.Error(w, r, http.StatusBadRequest, errors.New("Dict parameter not set"))
	}

	suggests, err := h.ucase.GetSuggest(query, updateService, dict)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetSuggest")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, suggests)
}
