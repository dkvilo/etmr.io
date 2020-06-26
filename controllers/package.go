package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dkvilo/etmr.io/core"
)

// BaseController structure
type BaseController struct {
	App *core.App
	W http.ResponseWriter
	Context context.Context
	ContextCancelFunc context.CancelFunc
	Status int
	Type string
}

// SendJSON Response Method
func (c *BaseController) SendJSON(content interface {}) {
	c.Type = "application/json"
	c.W.Header().Set("Content-Type", c.Type)
	c.W.WriteHeader(c.Status)
	json.NewEncoder(c.W).Encode(content)
}

