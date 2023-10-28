package controllers

import (
	"errors"
	"net/http"

	"github.com/itsjoetree/forest-life/helpers"
	"github.com/itsjoetree/forest-life/services"
)

var profile services.Profile

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		helpers.ErrorJSON(w, errors.New("user_id is required"), http.StatusBadRequest)
		return
	}

	profile, err := profile.GetProfileByUserId(userId)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to get profile"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"profile": profile})
}
