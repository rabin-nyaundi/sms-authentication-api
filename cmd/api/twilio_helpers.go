package main

import (
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

// sendOTP intializes twilio client sends otp
func (app *application) sendOTP(phonenumber string) error {

	// Load environment variables from .env
	if err := godotenv.Load("twilio.env"); err != nil {
		app.logger.PrintError(err, nil)
		return err
	}

	app.logger.PrintInfo("Twilio sms env loaded successfully", nil)

	client := twilio.NewRestClient()

	params := verify.CreateVerificationParams{}
	params.SetTo(phonenumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification("VAb8b1b718df0b522d53e4759de792a7f4", &params)
	if err != nil {
		app.logger.PrintError(err, nil)
		return err
	}

	if resp.Status != nil {
		app.logger.PrintInfo("response body", map[string]string{
			"account_sid": *resp.AccountSid,
			"status":      *resp.Status,
			"sent to":     *resp.To,
		})
	}

	app.logger.PrintInfo(*resp.Status, nil)

	return nil
}

func (app *application) verifyOTP(phonenumber, otpcode string) error {
	
	// Load environment variables from .env
	if err := godotenv.Load("twilio.env"); err != nil {
		app.logger.PrintError(err, nil)
		return err
	}

	app.logger.PrintInfo("Twilio sms env loaded successfully", nil)
	
	client := twilio.NewRestClient()

	params := verify.CreateVerificationCheckParams{}
	params.SetTo(phonenumber)
	params.SetCode(otpcode)

	resp, err := client.VerifyV2.CreateVerificationCheck("VAb8b1b718df0b522d53e4759de792a7f4", &params)

	if err != nil {
		app.logger.PrintError(err, nil)
		return err
	}

	if resp.Status != nil {
		app.logger.PrintInfo("response body", map[string]string{
			"status": *resp.Status,
			"to":     *resp.To,
		})
		return nil
	}
	app.logger.PrintInfo(*resp.Status, nil)
	return nil
}
