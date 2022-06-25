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
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mhkarimi1383/goExpenseTracker/configuration"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"github.com/mhkarimi1383/goExpenseTracker/types"
)

var (
	baseUrl string
)

func init() {
	cfg, err := configuration.GetConfig()
	if err != nil {
		logger.Fatalf(true, "error in initializing configuration: %v", err)
	}
	baseUrl = cfg.BaseURL
}

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
		Path:     "/",
		Domain:   "tracker.karimi.dev",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil || strings.HasPrefix(baseUrl, "https://"),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func ExtractTokenData(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return nil, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	//ok && token.Valid
	isValid := ok

	if isValid {
		return claims, nil
	} else {
		return nil, err
	}
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
