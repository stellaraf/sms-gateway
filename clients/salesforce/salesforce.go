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

func AuthenticateRequest(req *http.Request) error {
	env, err := util.LoadEnv()
	if err != nil {
		return err
	}
	key := req.Header.Get("x-api-key")
	return util.ValidateKey(key, env.AuthChecksumSalesforce)
}

func ValidateSalesforceData(codeRequest *types.SalesforceCodeRequest) (*types.SalesforceCodeRequest, error) {
	var data *types.SalesforceCodeRequest
	err := validator.New().Struct(codeRequest)
	if err != nil {
		return nil, err
	}
	country := countries.ByName(codeRequest.Country)
	if country == countries.Unknown {
		err = fmt.Errorf("country '%s' is invalid or unknown", codeRequest.Country)
		return nil, err
	}
	p, err := phonenumbers.Parse(codeRequest.PhoneNumber, country.Alpha2())
	if err != nil {
		return nil, err
	}
	phoneNumber := phonenumbers.Format(p, phonenumbers.E164)
	data = codeRequest
	data.PhoneNumber = phoneNumber
	return data, nil
}

func HandleSalesforceCodeRequest(codeRequest *types.SalesforceCodeRequest) error {
	env, err := util.LoadEnv()
	if err != nil {
		return err
	}
	body := fmt.Sprintf("Your Stellar Support PIN is %s\nReply STOP to opt out.", codeRequest.Code)
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
		return err
	}
	resBody, err := json.Marshal(*res)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Writer().Write(resBody)
	return nil
}
