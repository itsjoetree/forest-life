package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/itsjoetree/forest-life/helpers"
	"github.com/itsjoetree/forest-life/services"
)

var post services.Post

func LikePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err = post.LikePost(id, sessionId)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("serverError"))
	}

	helpers.WriteJSON(w, http.StatusOK, nil)
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err = post.UnlikePost(id, sessionId)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("serverError"))
	}

	helpers.WriteJSON(w, http.StatusOK, nil)
}

// GET/posts?author_id={author_id}
func GetPosts(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")

	// Author ID is required
	if authorId == "" {
		helpers.ErrorJSON(w, errors.New("author_id is required"), http.StatusBadRequest)
		return
	}

	posts, err := post.GetPosts(authorId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to get posts"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"posts": posts})
}

// GET.posts/{id}
func GetPostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := post.GetPostById(id)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to find post"), http.StatusNotFound)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, post)
}

// POST/posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	var newPost services.Post
	err = json.NewDecoder(r.Body).Decode(&newPost)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid JSON"))
		return
	}

	postCreated, err := post.CreatePost(newPost, sessionId)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to create post"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, postCreated)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Invalid JSON"))
		return
	}

	postUpdated, err := post.UpdatePost(id, post, sessionId)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to update post"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, postUpdated)
}

// DELETE/posts/{id}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	sessionId, err := auth.GetSessionId(r)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err = post.DeletePost(id, sessionId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.ErrorJSON(w, errors.New("Unable to delete post"), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "Post deleted")
}
