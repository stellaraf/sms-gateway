package types

type TestData struct {
	PhoneNumber       string `json:"phoneNumber"`
	AuthKeySalesforce string `json:"authKeySalesforce"`
}

type Environment struct {
	TwilioAccountSID       string   `json:"twilioAccountSid"`
	TwilioAuthToken        string   `json:"twilioAuthToken"`
	TwilioPhoneNumber      string   `json:"twilioPhoneNumber"`
	AuthChecksumSalesforce string   `json:"authChecksumSalesforce"`
	TestData               TestData `json:"testData"`
}
