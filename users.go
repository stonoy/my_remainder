package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stonoy/my_remainder/internal/database"
)

func (cfg *apiConfig) registerUser(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	decoder := json.NewDecoder(r.Body)
	var reqObj struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&reqObj)
	if err != nil {
		replyWithError(fmt.Sprintf("Can not decode %v", err), 400, w)
		return
	}

	if reqObj.Name == "" || reqObj.Email == "" || len(reqObj.Password) < 6 {
		replyWithError("check the inputs", 400, w)
		return
	}

	hashPasword, err := hashFromPassword(reqObj.Password)
	if err != nil {
		replyWithError(fmt.Sprintf("Error in hashing password %v", err), 500, w)
		return
	}

	newUser, err := cfg.dbQ.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqObj.Name,
		Email:     reqObj.Email,
		Password:  hashPasword,
	})
	if err != nil {
		replyWithError(fmt.Sprintf("error in creating user %v", err), 500, w)
		return
	}

	// generate token
	token, err := generateToken(newUser, cfg.jwt_secret)
	if err != nil {
		replyWithError(fmt.Sprintf("error in generating token %v", err), 500, w)
		return
	}

	type respStruct struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	replyWithJson(w, respStruct{
		User:  dbToRespUser(newUser),
		Token: token,
	}, 201)
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	// decode the request body
	decoder := json.NewDecoder(r.Body)
	var reqObj struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&reqObj)
	if err != nil {
		replyWithError(fmt.Sprintf("Can not decode %v", err), 400, w)
		return
	}

	if reqObj.Email == "" || len(reqObj.Password) < 6 {
		replyWithError("check the inputs", 400, w)
		return
	}

	theUser, err := cfg.dbQ.GetUserByEmail(r.Context(), reqObj.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			replyWithError("No such user found", 400, w)
		} else {
			replyWithError("error in GetUserByID", 500, w)
		}
		return
	}

	hasPasswordMatched := compareHashPassword(reqObj.Password, theUser.Password)
	if !hasPasswordMatched {
		replyWithError("password not matched", 401, w)
		return
	}

	token, err := generateToken(theUser, cfg.jwt_secret)
	if err != nil {
		replyWithError("error in generateToken", 500, w)
		return
	}

	type respStruct struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	replyWithJson(w, respStruct{
		User:  dbToRespUser(theUser),
		Token: token,
	}, 200)
}
