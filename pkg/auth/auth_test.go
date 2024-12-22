package auth

import (
	"strings"
	"testing"
	"time"
)

const secret = "my-secret-secret"

func TestJwtTokenExpired(t *testing.T) {
	accessTokenTTL := time.Minute * 1
	m := NewAuthManager([]byte(secret), accessTokenTTL)

	tokenPair, err := m.GetTokenPair("admin@gmail.com")
	if err != nil {
		t.Errorf("failed to get token pair %v", err)
	}

	// wait for the token to expire.
	time.Sleep(time.Second * 70)

	err = m.ValidateJwtToken(tokenPair.AccessToken)
	if err == nil || !strings.Contains(err.Error(), "jwt token is expired") {
		t.Errorf("token expired error is expected")
	}
}

func TestSuccessfullAccessTokenValidation(t *testing.T) {

}
