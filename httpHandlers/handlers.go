// all http handlers are here
package httpHandlers

import (
	"fmt"
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/mhkarimi1383/goExpenseTracker/configuration"
	"github.com/mhkarimi1383/goExpenseTracker/httpServer"
	"github.com/mhkarimi1383/goExpenseTracker/logger"
	"github.com/mhkarimi1383/goExpenseTracker/types"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
)

var (
	// global variable to store api address
	apiAddress string
	// global variable to store metric address
	metricAddress string
	// global variable to store the info of the application
	information       types.ApplicationInformation
	openIDUsernameKey string
)

// store needed variables from configuration at first import
func init() {
	cfg, err := configuration.GetConfig()
	if err != nil {
		logger.Fatalf(true, "error in initializing configuration: %v", err)
	}
	apiAddress = cfg.APIAddress
	metricAddress = cfg.MetricAddress
	information.Title = cfg.ApplicationTitle
	information.Description = cfg.ApplicationDescription
	openIDUsernameKey = cfg.OpenIDUsernameKey
}

// healthz is a function that returns a health-check (usually needed for Kubernetes/Docker) message to check if the server is running
func healthz(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	name := parameters["name"]
	if r.Header.Get("Content-Type") == "application/json" {
		// jsonResponse := types.UntypedMap{
		// 	"name":    name,
		// 	"msg": "I am healthy",
		// }
		jsonResponse := types.HealthzResponse{
			Name:    name,
			Message: "I am healthy",
		}
		err := responseWriter(w, &jsonResponse, http.StatusOK)
		if err != nil {
			logger.Warnf(true, "error while sending response %v", err)
		}
		return
	} else {
		strResponse := fmt.Sprintf("Hello %s", name)
		err := responseWriter(w, &strResponse, http.StatusOK)
		if err != nil {
			logger.Warnf(true, "error while sending response %v", err)
		}
		return
	}
}

// handler for will call healthz function
func healthzHandler() http.Handler {
	return http.HandlerFunc(healthz)
}

// notFound is a function that returns a not found message when the requested path is not found
// only for logging purpose
func notFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

// calling notFound function when the requested path is not found
func notFoundHandler() http.Handler {
	return http.HandlerFunc(notFound)
}

// Main function to start the server based on configuration
func RunServer() {
	figure := figure.NewFigure(information.Title, "doom", true)
	figure.Print()
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	}) // create a new middleware for metrics
	router := mux.NewRouter()                                                  // initialize router
	router.StrictSlash(true)                                                   // enable strict slash (/)
	router.Handle("/healthz/{name}", httpServer.WithLogging(healthzHandler())) // healthz route
	router.Handle("/", httpServer.WithLogging(indexHandler()))
	router.Handle("/flat-remix.css", httpServer.WithLogging(cssHandler()))
	router.Handle("/login", httpServer.WithLogging(loginHandler()))
	router.Handle("/auth/openid/callback", httpServer.WithLogging(callbackHandler()))
	router.NotFoundHandler = httpServer.WithLogging(notFoundHandler()) // setting not found handler
	mrouter := middlewarestd.Handler("", mdlw, router)                 // init router for metrics
	go func() {
		logger.Infof(false, "starting metric server on %v", metricAddress)
		logger.Fatalf(true, "error in metric http server: %v", http.ListenAndServe(metricAddress, promhttp.Handler()))
	}() // start metric server
	handler := cors.Default().Handler(mrouter) // add cors to the router
	logger.Infof(false, "starting main server on %v", apiAddress)
	logger.Fatalf(true, "error in main http server: %v", http.ListenAndServe(apiAddress, handler)) // start main server
}
