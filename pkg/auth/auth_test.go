package auth

import (
	"strings"
	"testing"
	"time"
)

const secret = "my-secret-secret"
const email = "admin@gmail.com"

func TestJwtTokenExpired(t *testing.T) {
	accessTokenTTL := time.Minute * 1
	m := NewAuthManager([]byte(secret), accessTokenTTL)

	tokenPair, err := m.GetTokens(email)
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
	m := NewAuthManager([]byte(secret), time.Second*30) // 30 seconds TTL

	tokenPair, err := m.GetTokens(email)
	if err != nil {
		t.Errorf("failed to get token pair %v", err)
	}

	// wait 15 seconds
	time.Sleep(time.Second * 15)

	if err := m.ValidateJwtToken(tokenPair.AccessToken); err != nil {
		t.Error(err)
	}
}
