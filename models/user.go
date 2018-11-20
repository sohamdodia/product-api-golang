package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type UserModel struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	Email        string        `json:"email" bson:"email"`
	Password     string        `json:"password,omitempty" bson:"-"`
	PasswordHash string        `json:"-" bson:"passwordHash"`
	Salt         string        `json:"-" bson:"salt"`
	Token        string        `json:"token,omitempty" bson:"-"`
}

func (u *UserModel) SetSaltedPassword(password string) error {
	salt := uuid.New().String()
	passwordBytes := []byte(password + salt)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash[:])
	u.Salt = salt

	return nil
}

func (u *UserModel) ComparePassword(password string) error {
	incoming := []byte(password + u.Salt)
	existing := []byte(u.PasswordHash)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err
}
