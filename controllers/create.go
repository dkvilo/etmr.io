package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/dkvilo/etmr.io/core"
	"github.com/dkvilo/etmr.io/models"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/go-playground/validator.v9"
)

// Create - Generate Custom Url and write to the db
// METHOD: POST
func (c *APIController) Create(w http.ResponseWriter, r *http.Request) {
	c.BaseController.W, c.BaseController.Status = w, http.StatusAccepted

	type BodyData struct {
		Slug string `json:"slug" validate:"required"`
		Link string `json:"link" validate:"link"`
	}

	var content BodyData
	_ = json.NewDecoder(r.Body).Decode(&content)

	// Validate r.Body Content
	v := validator.New()
	_ = v.RegisterValidation("link", func(fl validator.FieldLevel) bool {
		var status bool
		status, c.Error = regexp.MatchString("^(?:http(s)?:\\/\\/)[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$", content.Link)
		return status
	})
	c.Error = v.Struct(content)
	if c.Error != nil {
		for ind, _ := range c.Error.(validator.ValidationErrors) {
			c.BaseController.Status = http.StatusBadRequest
			c.BaseController.SendJSON(core.Response{
				Success: false,
				Message: fmt.Sprintf(
					"%s wants to be [%s]",
					c.Error.(validator.ValidationErrors)[0].Field(),
					c.Error.(validator.ValidationErrors)[0].ActualTag(),
				),
			})
			// Send Error one by one
			// Current Response Struct doesn't support array data structures
			if ind == 0 {
				break
			}
		}
		return
	}

	c.BaseController.Context, c.BaseController.ContextCancelFunc =
		context.WithTimeout(context.Background(), 10 * time.Second)

	linksCollection := c.BaseController.App.MongoClient.Database(os.Getenv("MONGODB_NAME")).Collection("links")

	var data interface{}
	if err := linksCollection.FindOne(c.BaseController.Context, bson.D{{Key: "slug",  Value: content.Slug}}).Decode(&data); err == nil {
		c.BaseController.Status = http.StatusFound
		c.BaseController.SendJSON(core.Response {
			Success: false,
			Message: "slug is already in use try different one",
		})
		return
	}

	_, c.Error = linksCollection.InsertOne(c.BaseController.Context, models.LinksData {
		 Schema: bson.D {
			 { Key: "link", Value: content.Link },
			 { Key: "slug", Value: content.Slug },
		 },
	}.Schema)

	if c.Error != nil {
		 c.BaseController.Status = http.StatusBadRequest
		 c.BaseController.SendJSON(core.Response{
			 Success: false,
			 Message: c.Error.Error(),
		 })
		 return
	}

	c.BaseController.Status = http.StatusCreated
	c.BaseController.SendJSON(core.Response{
		Success: true,
		Message: "Your Link was Successfully Shortened!",
		Data: map[string]interface{} {
			"short_link": fmt.Sprintf("%s/%s", os.Getenv("APP_DOMAIN"), content.Slug),
			"original_link": content.Link,
		},
	})

}
