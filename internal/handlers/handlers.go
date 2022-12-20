package handlers

import (
	"auth/internal/hashing"
	"auth/internal/tokens"
	"auth/internal/adapters/store/userstore"
	"context"

	"net/http"
	"fmt"
)

func Login(w http.ResponseWriter, r *http.Request, usersDb userstore.User, salt string) {
	login, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
	}
	dbUser, err := usersDb.Get(context.Background(), login)
	dbPwd := dbUser.Password
	role := dbUser.Role

	if err != nil || !hashing.VerifyPassword(pwd, dbPwd, salt) {
		w.WriteHeader(http.StatusForbidden)
	} else {
		accessCookie := &http.Cookie{
			Name:   "access_token",
			Value:  tokens.AccessToken(login, role),
		}
		refreshCookie := &http.Cookie{
			Name:   "refresh_token",
			Value:  tokens.RefreshToken(login, role),
		}
		http.SetCookie(w, accessCookie)
		http.SetCookie(w, refreshCookie)
		w.WriteHeader(http.StatusOK)
	}
}

func Verify(w http.ResponseWriter, r *http.Request) {
	accessCookie, err := r.Cookie("access_token")

	if err != nil {
		fmt.Fprintln(w, "Failed find access cookie: %#v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	refreshCookie, err := r.Cookie("refresh_token")

	if err != nil {
		fmt.Fprintln(w, "Failed find refresh cookie: %#v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	isAccessValid, parsedLogin, parsedRole := tokens.ValidateToken(accessCookie.Value)
	isRefreshValid, _, _ := tokens.ValidateToken(refreshCookie.Value)
	if isAccessValid {
		w.Header().Set("Login", parsedLogin)
		w.Header().Set("Role", parsedRole)
		w.WriteHeader(http.StatusOK)
		return
	}
	if isRefreshValid {
		newAccessCookie := &http.Cookie{
			Name:   "access_token",
			Value:  tokens.AccessToken(parsedLogin, parsedRole),
		}
		newRefreshCookie := &http.Cookie{
			Name:   "refresh_token",
			Value:  tokens.RefreshToken(parsedLogin, parsedRole),
		}
		http.SetCookie(w, newAccessCookie)
		http.SetCookie(w, newRefreshCookie)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusForbidden)
}