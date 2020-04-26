package controllers

import (
	"encoding/json"
	"fmt"

	"net/http"

	"../helpers"
	"../models"
	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SignInResponse ... sign in as User
type SignInResponse struct {
	Token string        `json:"token" bson:"token"`
	ID    bson.ObjectId `json:"id" bson:"id"`
}

// SignInAsUser ... sign in as User
type SignInAsUser struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type (

	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController ... returns a instance of UserController structure
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// CreateUser ... creates a new User resource
func (ad UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	User := models.User{}
	//prase json  of body and attach to admoin struct
	json.NewDecoder(r.Body).Decode(&User)

	//create id
	User.ID = bson.NewObjectId()

	fmt.Println(User.Username)
	fmt.Println(User.Password)

	encryptedPassword, er := helpers.Encrypt(User.Password)

	if er != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("Something goes wrong"))
		return
	}

	User.Password = encryptedPassword
	// write struct of admni to DB
	ad.session.DB("dibs").C("users").Insert(User)

	// build response for user
	token := helpers.GenerateToken(User.ID, "user")
	Res := SignInResponse{}
	Res.ID = User.ID
	Res.Token = token
	//

	// convert struct to JSON
	output, _ := json.Marshal(Res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", output)

}

// Signin ... sign in as User
func (ad UserController) Signin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	fmt.Println("Start Sign in")
	signIn := SignInAsUser{}
	json.NewDecoder(r.Body).Decode(&signIn)

	//verify username
	isValid, userPassword, userID := models.IsUserExist(signIn.Username)
	if isValid == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("Not valid username"))
		return
	}

	err := helpers.Compare(userPassword, signIn.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte("Not valid password"))
		return
	}

	token := helpers.GenerateToken(userID, "user")
	Res := SignInResponse{}
	Res.ID = userID
	Res.Token = token
	output, _ := json.Marshal(Res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", output)
}

// GetUser ... GetUser data
func (ad UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	fmt.Println("Start Get user data")

	// get user id from header
	id := p.ByName("id")
	fmt.Println("Start Get from id  is ", id)

	user := models.GetUser(id)
	if user.Username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("User is not Found"))
		return
	}
	output, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", output)
}
