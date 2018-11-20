package config

import (
	"gopkg.in/mgo.v2"
)

func GetSession() *mgo.Database {
	session, err := mgo.Dial("mongodb://" + Constants.Mongo.URL)

	if err != nil {
		panic(err)
	}

	dbSession := session.DB(Constants.Mongo.DBName)

	return dbSession
}
