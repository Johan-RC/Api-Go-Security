package jwt

import (
	"testing"
	"time"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	secret := "test-secret"
	now := time.Now()
	tok, err := GenerateAccessToken("user-1", "a@b.c", "USER", secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken error: %v", err)
	}

	claims, err := ParseToken(tok, secret)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Errorf("UserID = %q, esperado user-1", claims.UserID)
	}
	if claims.Email != "a@b.c" {
		t.Errorf("Email = %q", claims.Email)
	}
	if claims.ActorType != "USER" {
		t.Errorf("ActorType = %q, esperado USER", claims.ActorType)
	}
	if exp := claims.ExpiresAt.Time; !exp.After(now) {
		t.Errorf("ExpiresAt (%v) no es futuro", exp)
	}
	if claims.Issuer != "api-jeussairel" {
		t.Errorf("Issuer = %q", claims.Issuer)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	tok, err := GenerateAccessToken("u1", "e", "", "secret-a", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseToken(tok, "secret-b"); err == nil {
		t.Error("esperaba error con secret incorrecto")
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := ParseToken("no-es-un-jwt", "secret"); err == nil {
		t.Error("esperaba error para token inválido")
	}
}

func TestParseTokenExpired(t *testing.T) {
	secret := "test-secret"
	tok, err := GenerateAccessToken("u1", "a", "", secret, -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseToken(tok, secret)
	if err != ErrTokenExpired {
		t.Errorf("esperaba ErrTokenExpired, got %v", err)
	}
}