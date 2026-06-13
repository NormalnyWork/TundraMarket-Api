package auth_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	appauth "tundraMarket/internal/application/auth"
	domainadmin "tundraMarket/internal/domain/admin"
	domainauth "tundraMarket/internal/domain/auth"
	authinfrastructure "tundraMarket/internal/infrastructure/auth"
)

const testPassword = "test-admin-password"

type adminRepository struct {
	admin *domainadmin.Admin
	err   error
}

func (r adminRepository) GetByLogin(context.Context, string) (*domainadmin.Admin, error) {
	return r.admin, r.err
}

type tokenIssuer struct {
	claims domainauth.TokenClaims
}

func (i *tokenIssuer) Issue(claims domainauth.TokenClaims) (string, error) {
	i.claims = claims
	return "admin-token", nil
}

func (i *tokenIssuer) Verify(string) (*domainauth.TokenClaims, error) {
	return nil, errors.New("not implemented")
}

func TestAuthAdmin(t *testing.T) {
	passwordHash := hashPassword(t, testPassword)
	tokens := &tokenIssuer{}
	uc := appauth.NewUseCase(
		nil,
		nil,
		adminRepository{admin: domainadmin.New(1, "admin", passwordHash)},
		tokens,
		authinfrastructure.NewPasswordVerifier(),
	)

	token, err := uc.AuthAdmin(context.Background(), "admin", testPassword)
	if err != nil {
		t.Fatalf("AuthAdmin() error = %v", err)
	}
	if token != "admin-token" {
		t.Fatalf("AuthAdmin() token = %q, want %q", token, "admin-token")
	}
	if tokens.claims.Role != domainauth.RoleAdmin {
		t.Fatalf("issued role = %q, want %q", tokens.claims.Role, domainauth.RoleAdmin)
	}
}

func TestAuthAdminRejectsInvalidPassword(t *testing.T) {
	passwordHash := hashPassword(t, testPassword)
	uc := appauth.NewUseCase(
		nil,
		nil,
		adminRepository{admin: domainadmin.New(1, "admin", passwordHash)},
		&tokenIssuer{},
		authinfrastructure.NewPasswordVerifier(),
	)

	_, err := uc.AuthAdmin(context.Background(), "admin", "wrong-password")
	if !errors.Is(err, appauth.ErrUnauthorized) {
		t.Fatalf("AuthAdmin() error = %v, want %v", err, appauth.ErrUnauthorized)
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}
