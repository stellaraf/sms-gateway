package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/stellaraf/sms-gateway/clients/salesforce"
	"github.com/stellaraf/sms-gateway/types"
	"github.com/stellaraf/sms-gateway/util"
)

func Handler(writer http.ResponseWriter, req *http.Request) {
	err := util.ValidateMethod(req, http.MethodPost)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusMethodNotAllowed)
		return
	}

	_, err = util.LoadEnv()

	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		log.Fatal("failed to load environment variables:", err.Error())
		return
	}

	err = salesforce.AuthenticateRequest(req)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: "authentication failed"}, http.StatusUnauthorized)
		key := req.Header.Get("x-api-key")
		log.Fatalf("Authentication failed with key '%s'", key)
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		log.Fatal("failed to read body")
		return
	}
	log.Writer().Write(body)
	var codeRequest *types.SalesforceCodeRequest
	err = json.Unmarshal(body, &codeRequest)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		log.Fatal("failed to marshal JSON to struct")
		return
	}
	err = validator.New().Struct(codeRequest)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusBadRequest)
		log.Fatal("failed to validate request:", err.Error())
		return
	}
	err = salesforce.HandleSalesforceCodeRequest(codeRequest)
	if err != nil {
		util.SendError(writer, &types.APIErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		log.Fatal("failed to handle request:", err.Error())
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte{})
}
