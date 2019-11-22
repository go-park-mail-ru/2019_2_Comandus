package locationhttp

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/general/respond"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type LocationHandler struct {
	LocationUsecase		location.Usecase
	sanitizer      		*bluemonday.Policy
	logger         		*zap.SugaredLogger
	sessionStore   		sessions.Store
}

func NewLocationHandler(m *mux.Router, u location.Usecase, sanitizer *bluemonday.Policy, logger *zap.SugaredLogger, sessionStore sessions.Store) {
	handler := &LocationHandler{
		LocationUsecase:	u,
		sanitizer:      	sanitizer,
		logger:         	logger,
		sessionStore:   	sessionStore,
	}

	m.HandleFunc("/country-list", handler.HandleGetCountries).Methods(http.MethodGet, http.MethodOptions)
	m.HandleFunc("/city-list/{id:[0-9]+}", handler.HandleGetCities).Methods(http.MethodGet, http.MethodOptions)
}

func (h *LocationHandler) HandleGetCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	countries, err := h.LocationUsecase.CountryList()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCountries<-LocationUcase.CountryList()")
		respond.Error(w, r, http.StatusInternalServerError, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, countries)
}

func (h *LocationHandler) HandleGetCities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCities<-Atoi(wrong id)")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	cities, err := h.LocationUsecase.CityListByCountry(int64(id))
	if err != nil {
		err = errors.Wrapf(err, "HandleGetCities<-LocationUcase.CityListByCountry()")
		respond.Error(w, r, http.StatusBadRequest, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, cities)
}