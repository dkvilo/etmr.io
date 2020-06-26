package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dkvilo/etmr.io/core"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Redirect - redirect to original url
func (c *APIController) Redirect(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	c.BaseController.W, c.BaseController.Status = w, http.StatusAccepted

	c.BaseController.Context, c.BaseController.ContextCancelFunc = context.WithTimeout(context.Background(), 10 * time.Second)

	linksCollection := c.BaseController.App.MongoClient.Database(os.Getenv("MONGODB_NAME")).Collection("links")

	var data struct {
		ID string `json:"_id"`
		Slug string `json:"slug"`
		Link string `json:"link"`
	}

	if c.Error = linksCollection.FindOne(c.BaseController.Context, bson.D{{Key: "slug",  Value: params["slug"]}}).Decode(&data); c.Error != nil {
		c.BaseController.Status = http.StatusNotFound
		c.BaseController.SendJSON(core.Response {
			Success: false,
			Message: fmt.Sprintf("The Slug `%s` was not found", params["slug"]),
		})
	}

	c.BaseController.Status = http.StatusMovedPermanently
	http.Redirect(c.W, r, fmt.Sprintf("%s?ref=%s/%s", data.Link, os.Getenv("APP_DOMAIN"), data.Slug), c.BaseController.Status)
}
