package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/opentibiabr/login-server/src/api/models"
)

func TestClientUpdateRequired(t *testing.T) {
	tests := []struct {
		name    string
		minimum string
		current string
		want    bool
	}{
		{"disabled", "", "1.0.3", false},
		{"older", "1.0.4", "1.0.3", true},
		{"missing_old_client", "1.0.4", "", true},
		{"malformed_client", "1.0.4", "unknown", true},
		{"same", "1.0.4", "1.0.4", false},
		{"newer_patch", "1.0.4", "1.0.5", false},
		{"newer_minor", "1.0.4", "1.1.0", false},
		{"older_major", "2.0.0", "1.9.9", true},
		{"invalid_minimum", "latest", "1.0.3", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := clientUpdateRequired(test.minimum, test.current); got != test.want {
				t.Fatalf("clientUpdateRequired(%q, %q) = %v, want %v", test.minimum, test.current, got, test.want)
			}
		})
	}
}

func TestLoginRejectsOldAppBeforeAuthentication(t *testing.T) {
	t.Setenv("MIN_CLIENT_APP_VERSION", "1.0.4")
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", (&Api{}).login)

	requestBody, err := json.Marshal(models.RequestPayload{
		Type: "login", Email: "user@example.com", Password: "password123", AppVersion: "1.0.3",
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var response models.LoginErrorPayload
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(response.ErrorMessage, "1.0.4") || !strings.Contains(response.ErrorMessage, "zealot.website") {
		t.Fatalf("update guidance missing: %q", response.ErrorMessage)
	}
}
