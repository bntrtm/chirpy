package auth

import (
	"testing"
)

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