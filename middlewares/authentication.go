package middlewares

import (
	"net/http"
	"strings"

	"github.com/sohamdodia/product-api-golang/helper"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/sohamdodia/product-api-golang/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Middleware struct {
	dbSession *mgo.Database
}

func NewMiddleware(dbSession *mgo.Database) *Middleware {
	return &Middleware{dbSession}
}

func (m Middleware) AuthenticationMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				type CustomClaims struct {
					ID string `json:"id"`
					jwt.StandardClaims
				}

				id, err := helper.VerifyToken(bearerToken[1])

				if err != nil {
					helper.Response(w, http.StatusForbidden, false, "You are not authorized!", nil, err)
					return
				}

				if !bson.IsObjectIdHex(id) {
					helper.Response(w, http.StatusForbidden, false, "You are not authorized!", nil, nil)
					return
				}

				newID := bson.ObjectIdHex(id)

				user := models.UserModel{}
				err = m.dbSession.C("users").Find(bson.M{"_id": newID}).One(&user)

				if err != nil {
					helper.Response(w, http.StatusForbidden, false, "You are not authorized!", nil, err)
					return
				}
				context.Set(r, "user", user)
				next(w, r, p)
				return
			}
		}
		helper.Response(w, http.StatusForbidden, false, "You are not authorized!", nil, nil)
		return
	}
}
