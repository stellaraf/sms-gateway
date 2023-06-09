package util

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/stellaraf/sms-gateway/types"
)

func findGoMod(start string) (dir string, err error) {
	err = filepath.Walk(start, func(path string, file fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(file.Name(), "go.mod") {
			dir = filepath.Dir(path)
			return nil
		}
		return nil
	})
	return
}

func findProjectRoot() (root string, err error) {
	start, err := os.Getwd()
	if err != nil {
		return
	}
	start, err = filepath.Abs(start)
	if err != nil {
		return
	}
	for {
		dir, err := findGoMod(start)
		if err != nil {
			return "", err
		}
		if dir == "" {
			start = filepath.Dir(start)
			continue
		} else {
			root, err = filepath.Abs(dir)
			if err != nil {
				return "", err
			}
			break
		}
	}
	return
}

func loadDotEnv() (err error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return
	}
	envFile := filepath.Join(projectRoot, ".env")
	if _, err := os.Stat(envFile); err == nil {
		log.Println("loading environment variables from .env")
		err = godotenv.Load(envFile)
		if err != nil {
			return err
		}
	}
	return
}

func checkRequired(pairs ...[]string) (err error) {
	for _, p := range pairs {
		k := p[0]
		v := p[1]
		if v == "" {
			err = fmt.Errorf("environment variable '%s' is missing", k)
			return
		}
	}
	return
}

func LoadEnv() (env types.Environment, err error) {
	vercel := os.Getenv("VERCEL")
	if vercel == "" {
		err = loadDotEnv()
		if err != nil {
			return
		}
	}
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	authChecksumSalesforce := os.Getenv("AUTH_CHECKSUM_SALESFORCE")
	err = checkRequired(
		[]string{"TWILIO_ACCOUNT_SID", twilioAccountSid},
		[]string{"TWILIO_AUTH_TOKEN", twilioAuthToken},
		[]string{"TWILIO_PHONE_NUMBER", twilioPhoneNumber},
		[]string{"AUTH_CHECKSUM_SALESFORCE", authChecksumSalesforce},
	)
	if err != nil {
		return
	}
	testDataRaw := os.Getenv("TEST_DATA")
	var testData types.TestData
	err = json.Unmarshal([]byte(testDataRaw), &testData)
	if err != nil {
		return
	}
	env = types.Environment{
		AuthChecksumSalesforce: authChecksumSalesforce,
		TwilioAccountSID:       twilioAccountSid,
		TwilioAuthToken:        twilioAuthToken,
		TwilioPhoneNumber:      twilioPhoneNumber,
		TestData:               testData,
	}
	return
}
