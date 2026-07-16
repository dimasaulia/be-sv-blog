package controllers

import "net/http"

type AuthController interface {
	CurrentUser(w http.ResponseWriter, r *http.Request)
}
