package contractHttp

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	user_contract "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/user-contract"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ContractHandler struct {
	ContractUsecase	user_contract.Usecase
	sanitizer		*bluemonday.Policy
	logger			*zap.SugaredLogger
	sessionStore	sessions.Store
}

func NewContractHandler(m *mux.Router, cs user_contract.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &ContractHandler{
		ContractUsecase:	cs,
		sanitizer:			sanitizer,
		logger:				logger,
		sessionStore:		sessionStore,
	}

	m.HandleFunc("/responses/{id:[0-9]+}/contract", handler.HandleCreateContract).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+/done}", handler.HandleTickContractAsDone).Methods(http.MethodPut, http.MethodOptions)
	m.HandleFunc("/contract/{id:[0-9]+}", handler.HandleReviewContract).Methods(http.MethodPut, http.MethodOptions)
}

func (h * ContractHandler) HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateContract<-Body.Close:")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	//TODO: parse start end time here
	/*contract := new(model.Contract)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(contract)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract:")
		general.Error(w, r, http.StatusBadRequest, err)
		return
	}*/

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleCreateContract: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)


	if err := h.ContractUsecase.CreateContract(u, responseId); err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-ContractUsecase.CreateContract(): ")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: other responses should be denied
	respond.Respond(w, r, http.StatusOK, struct{}{})
}

func (h * ContractHandler) HandleTickContractAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleTickContractAsDone: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)
	if err := h.ContractUsecase.SetAsDone(u, contractId); err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-ContractUsecase.SetAsDone(): ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w,r, http.StatusOK, struct{}{})
}

func (h *ContractHandler) HandleReviewContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u, ok := r.Context().Value(respond.CtxKeyUser).(*model.User)
	if !ok {
		err := errors.Wrapf(errors.New("no user in context"),"HandleTickContractAsDone: ")
		respond.Error(w, r, http.StatusUnauthorized, err)
		return
	}


	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleReviewContract: ")
			respond.Error(w, r, http.StatusInternalServerError, err)
		}
	}()

	type Input struct {
		Grade int `json:"grade"`
	}

	decoder := json.NewDecoder(r.Body)
	input := new(Input)
	if err := decoder.Decode(input); err != nil {
		err = errors.Wrapf(err, "HandleReviewContract: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-strconv.Atoi: ")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)

	if err := h.ContractUsecase.ReviewContract(u, contractId, input.Grade); err != nil {
		err = errors.Wrapf(err, "HandleReviewContract<-contractUsecase.ReviewContract(): ")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w,r, http.StatusOK, struct{}{})
}

