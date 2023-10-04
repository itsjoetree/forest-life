package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/itsjoetree/forest-life/helpers"
	"github.com/itsjoetree/forest-life/services"
)

var post services.Post

// GET/posts?author_id={author_id}
func GetPosts(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")

	// Author ID is required
	if authorId == "" {
		helpers.MessageLogs.ErrorLog.Println("author_id is required")
		return
	}

	posts, err := post.GetPosts(authorId)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"posts": posts})
}

// POST/posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost services.Post
	err := json.NewDecoder(r.Body).Decode(&newPost)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	postCreated, err := post.CreatePost(newPost)

	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, postCreated)
}
