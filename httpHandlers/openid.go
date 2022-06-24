package httpHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/mhkarimi1383/goExpenseTracker/configuration"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"golang.org/x/oauth2"
)

var (
	oauthConfig *oauth2.Config
	provider    *oidc.Provider
)

func init() {
	cfg, err := configuration.GetConfig()
	if err != nil {
		logger.Fatalf(true, "error in initializing configuration: %v", err)
	}
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.OpenIDAddress)
	fmt.Printf("provider: %+v\n", provider)
	if err != nil {
		logger.Fatalf(true, "error in initializing openid provider: %v", err)
	}
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.OpenIDClientID,
		ClientSecret: cfg.OpenIDClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.BaseURL + "/auth/openid/callback",
		Scopes:       []string{oidc.ScopeOpenID, "email"},
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	state, err := randString(16)
	if err != nil {
		logger.Warnf(true, "error while generating state: %v", err)
		resp := http.StatusText(http.StatusInternalServerError)
		responseWriter(w, &resp, http.StatusInternalServerError)
	}
	setCallbackCookie(w, r, "state", state)
	http.Redirect(w, r, oauthConfig.AuthCodeURL(state), http.StatusFound)
}

func loginHandler() http.Handler {
	return http.HandlerFunc(login)
}

func callback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	if err != nil {
		resp := http.StatusText(http.StatusBadRequest) + ": " + "state cookie not set"
		responseWriter(w, &resp, http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		resp := http.StatusText(http.StatusBadRequest) + ": " + "state cookie not match"
		responseWriter(w, &resp, http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	oauth2Token, err := oauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		logger.Warnf(true, "failed to exchange with provider: %v", err)
		resp := http.StatusText(http.StatusInternalServerError) + ": " + "failed to exchange with provider"
		responseWriter(w, &resp, http.StatusInternalServerError)
		return
	}

	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		logger.Warnf(true, "failed to get user info: %v", err)
		resp := http.StatusText(http.StatusInternalServerError) + ": " + "failed to exchange with provider"
		responseWriter(w, &resp, http.StatusInternalServerError)
		return
	}

	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func callbackHandler() http.Handler {
	return http.HandlerFunc(callback)
}
