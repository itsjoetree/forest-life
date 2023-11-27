package controllers

import (
	"errors"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/itsjoetree/forest-life/helpers"
	"github.com/itsjoetree/forest-life/services"
)

var user services.User

func Follow(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.ErrorJSON(w, errors.New("idRequired"), http.StatusBadRequest)
		return
	}

	status, err := user.Follow(id, sessionId)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, status)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, nil)
}

func Unfollow(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.ErrorJSON(w, errors.New("idRequired"), http.StatusBadRequest)
		return
	}

	status, err := user.Unfollow(id, sessionId)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, err, status)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, nil)
}
