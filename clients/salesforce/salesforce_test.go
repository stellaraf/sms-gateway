package salesforce_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stellaraf/sms-gateway/clients/salesforce"
	"github.com/stellaraf/sms-gateway/types"
	"github.com/stellaraf/sms-gateway/util"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateSalesforceData(t *testing.T) {
	t.Run("salesforce validation succeeds", func(t *testing.T) {
		dataIn := &types.SalesforceCodeRequest{
			Country:     "United States",
			Code:        "1234",
			CaseNumber:  "123456",
			PhoneNumber: "+12223334444",
		}
		res, err := salesforce.ValidateSalesforceData(dataIn)
		assert.NoError(t, err)
		assert.Equal(t, dataIn.PhoneNumber, res.PhoneNumber)
	})
	t.Run("phone number conversion US", func(t *testing.T) {
		dataIn := &types.SalesforceCodeRequest{
			Country:     "United States",
			Code:        "1234",
			CaseNumber:  "123456",
			PhoneNumber: "(111) 222-3333",
		}
		expected := "+11112223333"
		res, err := salesforce.ValidateSalesforceData(dataIn)
		assert.NoError(t, err)
		assert.Equal(t, expected, res.PhoneNumber)
	})
	t.Run("phone number conversion UK", func(t *testing.T) {
		dataIn := &types.SalesforceCodeRequest{
			Country:     "United Kingdom",
			Code:        "1234",
			CaseNumber:  "123456",
			PhoneNumber: "+44 020 7946 0000",
		}
		expected := "+442079460000"
		res, err := salesforce.ValidateSalesforceData(dataIn)
		assert.NoError(t, err)
		assert.Equal(t, expected, res.PhoneNumber)
	})
}

func Test_AuthenticateRequest(t *testing.T) {
	t.Run("salesforce authentication passes", func(t *testing.T) {
		env, err := util.LoadEnv()
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "http://localhost", strings.NewReader(""))
		req.Header.Add("x-api-key", env.TestData.AuthKeySalesforce)
		assert.NoError(t, err)

		err = salesforce.AuthenticateRequest(req)
		assert.NoError(t, err)
	})

	t.Run("salesforce authentication fails", func(t *testing.T) {
		req, err := http.NewRequest("POST", "http://localhost", strings.NewReader(""))
		req.Header.Add("x-api-key", "wrong key")
		assert.NoError(t, err)

		err = salesforce.AuthenticateRequest(req)
		assert.Error(t, err)
	})
}
