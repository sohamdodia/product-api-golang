package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/context"

	"github.com/julienschmidt/httprouter"
	"github.com/sohamdodia/product-api-golang/helper"
	"github.com/sohamdodia/product-api-golang/models"
	"gopkg.in/mgo.v2"
)

type ProductController struct {
	dbSession *mgo.Database
}

func NewProductController(dbSession *mgo.Database) *ProductController {
	return &ProductController{dbSession}
}

func (pc ProductController) CreateProduct() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		type ValidationModel struct {
			Name     string `json:"name" valid:"required~Product name is required"`
			OldPrice int64  `json:"oldPrice" valid:"int,required"`
			NewPrice int64  `json:"newPrice" valid:"int,required"`
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		validationModelReader := bytes.NewReader(body)

		validationModel := ValidationModel{}

		json.NewDecoder(validationModelReader).Decode(&validationModel)

		_, err = govalidator.ValidateStruct(validationModel)

		if err != nil {
			helper.Response(w, http.StatusBadRequest, false, helper.TextTransform(err.Error()), nil, err)
			return
		}

		productModelReader := bytes.NewReader(body)
		product := models.ProductModel{}
		json.NewDecoder(productModelReader).Decode(&product)

		user := context.Get(r, "user").(models.UserModel)

		product.UserID = user.ID
		product.ID = bson.NewObjectId()
		fmt.Println(product)

		err = pc.dbSession.C("products").Insert(product)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		helper.Response(w, http.StatusOK, true, "Product added successfully.", product, nil)
	}
}

func (pc ProductController) DeleteProduct() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		id := p.ByName("id")

		if !bson.IsObjectIdHex(id) {
			helper.Response(w, http.StatusOK, false, "Enter a valid ID.", nil, nil)
			return
		}

		hexID := bson.ObjectIdHex(id)

		product := models.ProductModel{}
		user := context.Get(r, "user").(models.UserModel)
		err := pc.dbSession.C("products").Find(bson.M{"_id": hexID, "user_id": user.ID}).One(&product)

		if err != nil {
			helper.Response(w, http.StatusOK, false, "Product not found.", nil, err)
			return
		}

		err = pc.dbSession.C("products").Remove(bson.M{"_id": hexID})

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		helper.Response(w, http.StatusOK, true, "Product deleted successfully", nil, nil)
	}
}

func (pc ProductController) GetAllProducts() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		products := []models.ProductModel{}

		err := pc.dbSession.C("products").Pipe([]bson.M{
			{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "user_id",
					"foreignField": "_id",
					"as":           "user",
				},
			},
		}).All(&products)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		helper.Response(w, http.StatusOK, true, "Products fetched successfully", products, nil)
	}
}

func (pc ProductController) GetProduct() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		product := models.ProductModel{}

		id := p.ByName("id")

		if !bson.IsObjectIdHex(id) {
			helper.Response(w, http.StatusOK, false, "Enter a valid ID.", nil, nil)
			return
		}

		hexID := bson.ObjectIdHex(id)

		err := pc.dbSession.C("products").Pipe([]bson.M{
			{
				"$match": bson.M{
					"_id": hexID,
				},
			},
			{
				"$lookup": bson.M{
					"from":         "users",
					"localField":   "user_id",
					"foreignField": "_id",
					"as":           "user",
				},
			},
		}).One(&product)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		helper.Response(w, http.StatusOK, true, "Product fetched successfully!", product, nil)
	}
}

func (pc ProductController) UpdateProduct() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		type ValidationModel struct {
			Name     string `json:"name"`
			OldPrice int64  `json:"oldPrice" valid:"int"`
			NewPrice int64  `json:"newPrice" valid:"int"`
		}
		id := p.ByName("id")

		if !bson.IsObjectIdHex(id) {
			helper.Response(w, http.StatusOK, false, "Enter a valid ID.", nil, nil)
			return
		}
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		validationModelReader := bytes.NewReader(body)

		validationModel := ValidationModel{}

		json.NewDecoder(validationModelReader).Decode(&validationModel)

		_, err = govalidator.ValidateStruct(validationModel)

		if err != nil {
			helper.Response(w, http.StatusBadRequest, false, helper.TextTransform(err.Error()), nil, err)
			return
		}

		productModelReader := bytes.NewReader(body)
		product := models.ProductModel{}
		json.NewDecoder(productModelReader).Decode(&product)

		hexID := bson.ObjectIdHex(id)
		err = pc.dbSession.C("products").Find(bson.M{"_id": hexID}).One(&product)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		user := context.Get(r, "user").(models.UserModel)
		count, err := pc.dbSession.C("products").Find(bson.M{"_id": hexID, "user_id": user.ID}).Count()

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		if count == 0 {
			helper.Response(w, http.StatusNotFound, false, "Product not found!", nil, err)
			return
		}

		err = pc.dbSession.C("products").Update(bson.M{"_id": hexID}, &product)

		if err != nil {
			helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
			return
		}

		helper.Response(w, http.StatusOK, true, "Product updated successfully!", product, nil)
	}
}
