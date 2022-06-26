package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func SendJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)

	w.Write(buf.Bytes())
}

// JSON is a map alias
type JSON map[string]interface{}

// SendErrorJSON sends error as json
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string) {
	zap.L().Warn(details, zap.Error(err), zap.Int("httpStatusCode", httpStatusCode), zap.String("url", r.URL.String()), zap.String("method", r.Method))
	SendJSON(w, r, httpStatusCode, JSON{"error": err.Error(), "details": details})
}

func SendNoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
