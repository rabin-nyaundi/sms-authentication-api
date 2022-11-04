package data

type OTPData struct {
	PhoneNumber string `json:"phonenumber,omitempty"`
}

type VerifyOTP struct {
	OtpData OTPData `json:"otpdata,omitempty"`
	OtpCode string  `json:"otp,omitempty"`
}
