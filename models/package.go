package models

import "go.mongodb.org/mongo-driver/bson"

// BaseModel type
type BaseModel struct {
	Data bson.D
}
