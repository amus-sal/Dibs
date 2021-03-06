package controllers

import (
	"encoding/json"
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"net/http"

	"../helpers"
	"../models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

// SignIn ... sign in as Admin
type SignIn struct {
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
	Password string `json:"password" bson:"password"`
}

type (

	// AdminController represents the controller for operating on the admin resource
	AdminController struct {
		session *mgo.Session
	}
)

// NewAdminController ... returns a instance of UserController structure
func NewAdminController(session *mgo.Session) *AdminController {
	return &AdminController{session}
}

// CreateAdmin ... creates a new admin resource
func (ad AdminController) CreateAdmin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// tokn := r.Header.Get("Authorization")

	// _, err := helpers.VerifyToken(tokn)

	// if err != nil {
	// 	w.WriteHeader(404)
	// 	w.Write([]byte("Missing token"))
	// } else {
	admin := models.Admin{}
	//prase json  of body and attach to admoin struct
	json.NewDecoder(r.Body).Decode(&admin)

	//create id
	admin.ID = bson.NewObjectId()

	// write struct of admni to DB
	ad.session.DB("dibs").C("admins").Insert(admin)

	// convert struct to JSON
	res := helpers.ResController{Res: w}
	res.SendResponse("the admin is created", 200)

	// fmt.Fprintf(w, "%s", uj)
}

// Signin ... sign in as Admin
func (ad AdminController) Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println("Start Sign in")
	signIn := SignIn{}
	json.NewDecoder(r.Body).Decode(&signIn)

	isAdmin, password, id, role := models.IsAdminExist(signIn.Username, ad.session)
	if isAdmin != true {
		if signIn.Username == "admin" && signIn.Password == "admin" {
			// create user first time

			admin := models.Admin{}
			admin.Username = signIn.Username
			admin.Password = signIn.Password
			admin.Role = "admin"
			encryptedPassword, er := helpers.Encrypt(admin.Password)

			if er != nil {
				res := helpers.ResController{Res: w}
				res.SendResponse("Something went wrong", 404)
			} else {
				//create id
				admin.ID = bson.NewObjectId()
				admin.Password = encryptedPassword
				// write struct of admin to DB
				ad.session.DB("dibs").C("admins").Insert(admin)

				token := helpers.GenerateToken(admin.ID, admin.Role)
				fmt.Println("Token Generated is", token)

				Res := SignInResponse{}
				Res.ID = admin.ID
				Res.Token = token
				output, _ := json.Marshal(Res)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				fmt.Fprintf(w, "%s", output)
			}
		} else {
			res := helpers.ResController{Res: w}
			res.SendResponse("Not Authorized", 404)
		}
	} else {
		err := helpers.Compare(password, signIn.Password)
		if err != nil {
			res := helpers.ResController{Res: w}
			res.SendResponse("Not Authorized", 404)
		} else {
			token := helpers.GenerateToken(id, role)
			Res := SignInResponse{}
			Res.ID = id
			Res.Token = token
			output, _ := json.Marshal(Res)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			fmt.Fprintf(w, "%s", output)

		}
	}
}
