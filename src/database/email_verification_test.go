package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/opentibiabr/login-server/src/serviceerrors"
	"testing"
)

func TestEmailVerificationOnAuthentication(t *testing.T) {
	for _, tc := range []struct {
		name, email string
		verified    int
		allowed     bool
	}{
		{"verified", "player@example.com", 1, true},
		{"unverified", "player@example.com", 0, false},
		{"gm exception", "admin@zealot.website", 0, true},
		{"name cannot impersonate gm", "player@example.com", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			mock.ExpectQuery("SELECT id, type, premdays, lastday, email, email_verified FROM accounts").WillReturnRows(sqlmock.NewRows([]string{"id", "type", "premdays", "lastday", "email", "email_verified"}).AddRow(1, 1, 0, 0, tc.email, tc.verified))
			acc := &Account{Email: "admin@zealot.website", Password: "test"}
			err = acc.Authenticate(db)
			if tc.allowed && err != nil {
				t.Fatal(err)
			}
			if !tc.allowed {
				e, ok := serviceerrors.FromError(err)
				if !ok || e.Name != "EMAIL_NOT_VERIFIED" {
					t.Fatalf("expected verification error, got %v", err)
				}
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
