package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// HASH TESTS

func WasHashed(t *testing.T) {
	// passes if hashed password is indeed different from original password
	password := "cheetohDeadbolt123"
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	if password == hashedPass {
		t.Error("password was not hashed")
	}
}

func TestHashUnequal(t *testing.T) {
	// passes if CheckPasswordHash returns not nil as expected
	password := "cheetohDeadbolt123"
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	altPassword := "cheetohDeadbolt124"
	err = CheckPasswordHash(altPassword, hashedPass)
	if err == nil {
		t.Error("password should not have matched, but did")
	}
}

func TestHashEqual(t *testing.T) {
	// passes if CheckPasswordHash returns nil as expected
	password := "cheetohDeadbolt123"
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	err = CheckPasswordHash(password, hashedPass)
	if err != nil {
		t.Error("password should have matched, but did not")
	}

}

// JWT TESTS

func JWTRejectExpired(t *testing.T) {
	// passes if an expired JWT is properly rejected
	userID := uuid.New()
	tokenSecret := "very-secret-secret"
	expiration := time.Second * 2
	token, err := MakeJWT(userID, tokenSecret, expiration)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	_, err = ValidateJWT(token, "very-secret-secret")
	if err == nil {
		t.Error("expired JWT not rejected")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
