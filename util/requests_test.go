package util_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stellaraf/sms-gateway/util"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateMethod(t *testing.T) {
	t.Run("method validation passes", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://localhost", strings.NewReader(""))
		assert.NoError(t, err)
		err = util.ValidateMethod(req, "POST")
		assert.NoError(t, err)
	})
	t.Run("method validation fails", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost", strings.NewReader(""))
		assert.NoError(t, err)
		err = util.ValidateMethod(req, "POST")
		assert.Error(t, err)
	})
}

func Test_ValidateKey(t *testing.T) {
	t.Run("key validation passes", func(t *testing.T) {
		env, err := util.LoadEnv()
		assert.NoError(t, err)
		err = util.ValidateKey(env.TestData.AuthKeySalesforce, env.AuthChecksumSalesforce)
		assert.NoError(t, err)
	})
	t.Run("key validation fails", func(t *testing.T) {
		env, err := util.LoadEnv()
		assert.NoError(t, err)
		err = util.ValidateKey("wrong key", env.AuthChecksumSalesforce)
		assert.Error(t, err)
	})
}
