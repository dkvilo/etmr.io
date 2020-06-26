package controller

import (
	"net/http"
)

// APIController base
type APIController struct {
	BaseController
	Error error
}

// URL - Generate Custom Url and write to the db
// METHOD: MIXED
func (c *APIController) URL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.Create(w, r)
	}
}

// RedirectURL - redirect to original url
func (c *APIController) RedirectURL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.Redirect(w, r)
	}
}


