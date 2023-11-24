package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/itsjoetree/forest-life/helpers"
	"github.com/itsjoetree/forest-life/services"
)

var auth services.Auth

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds services.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	cookie, status, err := auth.SignIn(creds)

	if err != nil {
		helpers.ErrorJSON(w, err, status)
		return
	}

	http.SetCookie(w, cookie)
	helpers.WriteJSON(w, http.StatusOK, nil)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var payload services.Auth

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	cookie, status, err := auth.SignUp(payload)

	if err != nil {
		helpers.ErrorJSON(w, err, status)
		return
	}

	http.SetCookie(w, cookie)
	helpers.WriteJSON(w, http.StatusOK, nil)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		helpers.ErrorJSON(w, errors.New("bad request"), http.StatusBadRequest)

		return
	}

	cookie, status, err := auth.Refresh(c.Value)

	if err != nil {
		helpers.ErrorJSON(w, err, status)
		return
	}

	http.SetCookie(w, cookie)
	helpers.WriteJSON(w, http.StatusOK, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		helpers.ErrorJSON(w, errors.New("bad request"), http.StatusBadRequest)
		return
	}

	cookie, status, err := auth.Logout(c.Value)

	if err != nil {
		helpers.ErrorJSON(w, err, status)
		return
	}

	// Sets an empty cookie in session
	http.SetCookie(w, cookie)
	helpers.WriteJSON(w, http.StatusOK, nil)
}
