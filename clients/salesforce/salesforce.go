package salesforce

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/biter777/countries"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
	"github.com/stellaraf/sms-gateway/types"
	"github.com/stellaraf/sms-gateway/util"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func AuthenticateRequest(req *http.Request) (err error) {
	env, err := util.LoadEnv()
	if err != nil {
		return
	}
	key := req.Header.Get("x-api-key")
	err = util.ValidateKey(key, env.AuthChecksumSalesforce)
	return
}

func ValidateSalesforceData(codeRequest *types.SalesforceCodeRequest) (data *types.SalesforceCodeRequest, err error) {
	err = validator.New().Struct(codeRequest)
	if err != nil {
		return
	}
	country := countries.ByName(codeRequest.Country)
	if country == countries.Unknown {
		err = fmt.Errorf("country '%s' is invalid or unknown", codeRequest.Country)
		return
	}
	p, err := phonenumbers.Parse(codeRequest.PhoneNumber, country.Alpha2())
	phoneNumber := phonenumbers.Format(p, phonenumbers.E164)
	data = codeRequest
	data.PhoneNumber = phoneNumber
	return
}

func HandleSalesforceCodeRequest(codeRequest *types.SalesforceCodeRequest) (err error) {
	env, err := util.LoadEnv()
	if err != nil {
		return
	}
	body := fmt.Sprintf("Your Support PIN for Case %s is %s. Reply STOP to opt out.", codeRequest.CaseNumber, codeRequest.Code)
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSID,
		Password: env.TwilioAuthToken,
	})
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(codeRequest.PhoneNumber)
	params.SetFrom(env.TwilioPhoneNumber)
	params.SetBody(body)
	res, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		return
	}
	resBody, err := json.Marshal(*res)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Writer().Write(resBody)
	return
}
