package main

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"github.com/sohamdodia/product-api-golang/config"
	"github.com/sohamdodia/product-api-golang/controllers"
	"github.com/sohamdodia/product-api-golang/middlewares"
	mgo "gopkg.in/mgo.v2"
)

var dbSession *mgo.Database

func init() {
	dbSession = config.GetSession()
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(dbSession)
	pc := controllers.NewProductController(dbSession)
	m := middlewares.NewMiddleware(dbSession)

	r.POST("/signin", uc.Signin)
	r.POST("/signup", uc.Signup)
	r.GET("/product", m.AuthenticationMiddleware(pc.GetAllProducts()))
	r.GET("/product/:id", m.AuthenticationMiddleware(pc.GetProduct()))
	r.POST("/product", m.AuthenticationMiddleware(pc.CreateProduct()))
	r.PUT("/product/:id", m.AuthenticationMiddleware(pc.UpdateProduct()))
	r.DELETE("/product/:id", m.AuthenticationMiddleware(pc.DeleteProduct()))

	http.ListenAndServe(":"+strconv.Itoa(config.Constants.PORT), r)
}
