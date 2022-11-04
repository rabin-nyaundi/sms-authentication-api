package main

import (
	"fmt"
	"net/http"
)

type envelope map[string]interface{}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := envelope{
		"status":      "available",
		"environment": app.config.env,
		"version":     Version,
	}

	err := app.writeJSON(w, http.StatusOK, status)

	if err != nil {
		app.logger.PrintError(err, map[string]string{})
	}
}

func (app *application) sendSMSHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		PhoneNumber string `json:"phonenumber,omitempty"`
	}
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.sendOTP(input.PhoneNumber)

	if err != nil {
		app.logger.PrintError(err, nil)
	}

	fmt.Println(input.PhoneNumber)
}

func (app *application) verifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PhoneNumber string `json:"phonenumber,omitempty"`
		Otp         string `json:"otp,omitempty"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.logger.PrintError(err, nil)

	}
	err = app.verifyOTP(input.PhoneNumber, input.Otp)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
	app.logger.PrintInfo(fmt.Sprint(input), nil)
}
