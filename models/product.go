package models

import (
	"gopkg.in/mgo.v2/bson"
)

type ProductModel struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	OldPrice int64         `json:"oldPrice" bson:"old_price"`
	NewPrice int64         `json:"newPrice" bson:"new_price"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id"`
	User     []UserModel   `json:"user,omitempty"`
}
