package respond

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_Comandus/monitoring"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ctxKey int8

const (
	SessionName        = "user-session"
	CtxKeyUser              ctxKey = iota
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	// TODO: fix logger
	/*ctx := r.Context()
	reqID := ctx.Value("rIDKey").(string)
	s.logger.Infof("Request ID: %s | error : %s", reqID , err.Error())*/
	log.Println(err)
	Respond(w, r, code, map[string]string{"error": errors.Cause(err).Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	path := r.URL.String()
	i := strings.Index(path, "?")

	var cleanPath string
	if i > -1 {
		cleanPath = path[:i]
	} else {
		cleanPath = path
	}
	monitoring.Hits.WithLabelValues(strconv.Itoa(code), cleanPath).Inc()

	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
