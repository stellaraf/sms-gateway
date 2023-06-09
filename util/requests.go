package util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stellaraf/sms-gateway/types"
)

func ValidateMethod(request *http.Request, methods ...string) (err error) {
	method := strings.ToUpper(request.Method)
	for _, m := range methods {
		u := strings.ToUpper(m)
		if u == method {
			return
		}
	}
	err = fmt.Errorf("method '%s' is not allowed", method)
	return
}

func SendError(writer http.ResponseWriter, errorResponse *types.APIErrorResponse, status int) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(status)
	body, err := json.Marshal(errorResponse)
	if err != nil {
		errorResponse = &types.APIErrorResponse{
			Error: err.Error(),
		}
		body, err = json.Marshal(errorResponse)
		if err != nil {
			writer.Write([]byte{})
			log.Fatal(err)
		}
	}
	writer.Write(body)
}

func ValidateKey(key string, checksum string) (err error) {
	hashB := sha256.New()
	reader := strings.NewReader(key)
	_, err = io.Copy(hashB, reader)
	if err != nil {
		return
	}
	hash := fmt.Sprintf("%x", hashB.Sum(nil))
	if hash != checksum {
		err = fmt.Errorf("authentication failed")
		return
	}
	return
}
