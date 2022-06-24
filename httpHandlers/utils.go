package httpHandlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/mhkarimi1383/goExpenseTracker/types"
)

// responseWriter is a function to send response to the client with the given status code
// and decide whether to send the response as json or string then set the content type
func responseWriter[R string | types.UntypedMap | types.Response](w http.ResponseWriter, response *R, status int) error {
	w.WriteHeader(status)
	if reflect.ValueOf(*response).Kind() == reflect.Struct || reflect.ValueOf(*response).Kind() == reflect.Map {
		err := json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			return errors.New("error in encoding response")
		}
	} else {
		_, err := fmt.Fprintf(w, "%s", *response)
		if err != nil {
			return errors.New("error in writing response")
		}
	}
	return nil
}

func operatorCheckboxTranslator(v string) string {
	if v == "on" {
		return "+"
	}
	return "-"
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
