package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"

	"auth/internal/handlers"
)

func TestLoginTestifySuccess(t *testing.T) {
	testCases := []struct {
		NameTest string
		Email string
		Password string
	}{
		{
			NameTest: "successful user",
			Email: "julia",
			Password: "1234",
		},
	}

	handler := http.HandlerFunc(handlers.Login)
	for _, testCase := range testCases {
		t.Run(testCase.NameTest, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/api/v1/login", nil)

			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestLoginTestifyFailed(t *testing.T) {
	testCases := []struct {
		NameTest string
		Email string
		Password string
		StatusCode int
	}{
		{
			NameTest: "wrong login",
			Email: "login",
			Password: "password",
			StatusCode: http.StatusForbidden,
		},
		{
			NameTest: "wrong password",
			Email: "julia",
			Password: "password",
			StatusCode: http.StatusForbidden,
		},
	}

	handler := http.HandlerFunc(handlers.Login)
	for _, testCase := range testCases {
		t.Run(testCase.NameTest, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/api/v1/login", nil)

			handler.ServeHTTP(rec, req)
			assert.Equal(t, testCase.StatusCode, rec.Code)
		})
	}
}

func TestVerifyTestifySuccess(t *testing.T) {
	testCases := []struct {
		NameTest string
		Cookies []*http.Cookie
	}{
		{
			NameTest: "successful verify",
			Cookies: []*http.Cookie{
				&http.Cookie{
					Name: "access_token",
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpbiI6Imp1bGlhIiwiUm9sZSI6ImFkbWluaXN0cmF0b3IiLCJleHAiOjE2Njc1OTk5OTB9.SHxAdS3YktriWeVyyeWtGqCYZ8s-BBJwHnsd4SOpanI",
				},
				&http.Cookie{
					Name: "refresh_token",
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpbiI6Imp1bGlhIiwiUm9sZSI6ImFkbWluaXN0cmF0b3IiLCJleHAiOjE2Njc2MDM1MzB9.9UojCXPEh8YHotcdiq5zGni09Fy83uaVvRUafkggbQ0",
				},
			},
		},	
	}

	handler := http.HandlerFunc(handlers.Verify)
	for _, testCase := range testCases {
		t.Run(testCase.NameTest, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/api/v1/verify", nil)

			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestVerifyTestifyFailed(t *testing.T) {
	testCases := []struct {
		NameTest string
		Cookies []*http.Cookie
		statusCode int
	}{
		{
			NameTest: "no refresh cookie",
			Cookies: []*http.Cookie{
				&http.Cookie{
					Name: "access_token",
					Value: "any_value",
				},
			},
			statusCode: http.StatusNotFound,
		},	
		{
			NameTest: "no access cookie",
			Cookies: []*http.Cookie{
				&http.Cookie{
					Name: "refresh_token",
					Value: "any_value",
				},
			},
			statusCode: http.StatusNotFound,
		},	
	}

	handler := http.HandlerFunc(handlers.Verify)
	for _, testCase := range testCases {
		t.Run(testCase.NameTest, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/api/v1/verify", nil)
			for _, c := range testCase.Cookies {
				req.AddCookie(c)
			}
			handler.ServeHTTP(rec, req)
			assert.Equal(t, testCase.statusCode, rec.Code)
		})
	}	
}