package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/health", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/send-otp", app.sendSMSHandler)
	router.HandlerFunc(http.MethodPost, "/api/verify-otp", app.verifyOTPHandler)

	return router
}
