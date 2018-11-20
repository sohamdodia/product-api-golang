package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"github.com/sohamdodia/product-api-golang/helper"
	"github.com/sohamdodia/product-api-golang/models"
	"gopkg.in/mgo.v2"
)

type UserController struct {
	dbSession *mgo.Database
}

func NewUserController(dbSession *mgo.Database) *UserController {
	return &UserController{dbSession}
}

func (uc UserController) Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	type ValidationModel struct {
		Email    string `json:"email" valid:"email~Enter a valid email.,required~Email is required."`
		Password string `json:"password" valid:"required~Password is required"`
	}

	validationModel := ValidationModel{}

	json.NewDecoder(r.Body).Decode(&validationModel)

	_, err := govalidator.ValidateStruct(validationModel)
	if err != nil {
		helper.Response(w, http.StatusBadRequest, false, helper.TextTransform(err.Error()), nil, err)
		return
	}

	user := models.UserModel{}

	count, err := uc.dbSession.C("users").Find(bson.M{"email": validationModel.Email}).Count()

	if err != nil {
		helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, err)
		return
	}

	if count == 0 {
		helper.Response(w, http.StatusBadRequest, false, "User not found.", nil, nil)
		return
	}

	err = uc.dbSession.C("users").Find(bson.M{"email": validationModel.Email}).One(&user)

	if err != nil {
		helper.Response(w, http.StatusInternalServerError, false, "Something went wrong!", nil, nil)
		return
	}

	err = user.ComparePassword(validationModel.Password)

	if err != nil {
		helper.Response(w, http.StatusUnprocessableEntity, false, "Invalid Credentials! Please check your credentials.", nil, err)
		return
	}

	token, err := helper.GenerateToken(user.ID.Hex())

	if err != nil {
		helper.Response(w, http.StatusInternalServerError, false, "Something went wront!", nil, err)
		return
	}

	user.Token = token

	helper.Response(w, http.StatusOK, true, "Logged in successfully", user, nil)

}

func (uc UserController) Signup(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	type ValidationModel struct {
		Name     string `json:"name" valid:"alpha,required~Name is required"`
		Email    string `json:"email" valid:"email~Enter a valid email.,required~Email is required."`
		Password string `json:"password" valid:"required~Password is required"`
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

	userModelReader := bytes.NewReader(body)
	user := models.UserModel{}
	json.NewDecoder(userModelReader).Decode(&user)

	count, err := uc.dbSession.C("users").Find(bson.M{"email": validationModel.Email}).Count()

	if count != 0 {
		helper.Response(w, http.StatusBadRequest, false, "Email already exist!", nil, nil)
		return
	}

	uc.dbSession.C("users").Find(bson.M{"email": validationModel.Email}).One(&user)

	user.SetSaltedPassword(validationModel.Password)
	user.ID = bson.NewObjectId()

	uc.dbSession.C("users").Insert(user)

	token, err := helper.GenerateToken(user.ID.Hex())

	if err != nil {
		helper.Response(w, http.StatusInternalServerError, false, "Something went wront!", nil, err)
		return
	}

	user.Token = token

	helper.Response(w, http.StatusOK, true, "Signed up successfully!", user, nil)
}
