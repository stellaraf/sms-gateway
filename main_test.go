package main

import (
	"testing"

	"github.com/stellaraf/sms-gateway/clients/salesforce"
	"github.com/stellaraf/sms-gateway/types"
	"github.com/stellaraf/sms-gateway/util"
	"github.com/stretchr/testify/assert"
)

func Test_Salesforce(t *testing.T) {
	t.Run("send code", func(t *testing.T) {
		env, err := util.LoadEnv()
		assert.NoError(t, err)
		codeRequest := &types.SalesforceCodeRequest{
			CaseNumber:  "12345",
			Code:        "6789",
			PhoneNumber: env.TestData.PhoneNumber,
		}
		err = salesforce.HandleSalesforceCodeRequest(codeRequest)
		assert.NoError(t, err)
	})
}
