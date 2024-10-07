package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stonoy/my_remainder/internal/database"
)

func (cfg *apiConfig) createRemainders(w http.ResponseWriter, r *http.Request, user database.User) {
	type reqStruct struct {
		Subject      string `json:"subject"`
		Description  string `json:"description"`
		Has_Priority bool   `json:"has_priority"`
		Timing       string `json:"timing"`
	}

	decoder := json.NewDecoder(r.Body)
	reqObj := reqStruct{}
	err := decoder.Decode(&reqObj)
	if err != nil {
		replyWithError(fmt.Sprintf("Can not decode %v", err), 400, w)
		return
	}

	theTime, err := strToTime(reqObj.Timing)
	if err != nil {
		replyWithError(fmt.Sprintf("%v", err), 400, w)
		return
	}

	theRemainder, err := cfg.dbQ.CreateRemainder(r.Context(), database.CreateRemainderParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Subject:     reqObj.Subject,
		Description: reqObj.Description,
		Timing:      theTime,
		HasPriority: reqObj.Has_Priority,
		Userid:      user.ID,
	})
	if err != nil {
		replyWithError(fmt.Sprintf("err in CreateRemainder -> %v", err), 500, w)
		return
	}

	type respStruct struct {
		Success   bool      `json:"success"`
		Remainder Remainder `json:"remainder"`
	}

	replyWithJson(w, respStruct{
		Success:   true,
		Remainder: dbToRespRemainder(theRemainder),
	}, 201)

}

func (cfg *apiConfig) getRemaindersByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	remainders, err := cfg.dbQ.GetRemaindersByUser(r.Context(), user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			replyWithError("No remainders found", 200, w)
		} else {
			replyWithError(fmt.Sprintf("err in GetRemaindersByUser -> %v", err), 500, w)
		}
		return
	}

	type respStruct struct {
		Remainders []Remainder `json:"remainders"`
	}

	replyWithJson(w, respStruct{
		Remainders: dbToRespRemainders(remainders),
	}, 201)

}

func (cfg *apiConfig) getRemainderByID(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "ID")

	id, err := uuid.Parse(idStr)
	if err != nil {
		replyWithError("not a valid remainder id", 400, w)
		return
	}

	remainder, err := cfg.dbQ.GetRemainderByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			replyWithError("No remainder found", 200, w)
		} else {
			replyWithError(fmt.Sprintf("err in GetRemainderByID -> %v", err), 500, w)
		}
		return
	}

	type respStruct struct {
		Remainder Remainder `json:"remainder"`
	}

	replyWithJson(w, respStruct{
		Remainder: dbToRespRemainder(remainder),
	}, 200)
}

func (cfg *apiConfig) updateRemainder(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "ID")

	id, err := uuid.Parse(idStr)
	if err != nil {
		replyWithError("not a valid remainder id", 400, w)
	}

	type reqStruct struct {
		Subject      string `json:"subject"`
		Description  string `json:"description"`
		Has_Priority bool   `json:"has_priority"`
		Timing       string `json:"timing"`
	}

	decoder := json.NewDecoder(r.Body)
	reqObj := reqStruct{}
	err = decoder.Decode(&reqObj)
	if err != nil {
		replyWithError(fmt.Sprintf("Can not decode %v", err), 400, w)
		return
	}

	theTime, err := strToTime(reqObj.Timing)
	if err != nil {
		replyWithError(fmt.Sprintf("%v", err), 400, w)
		return
	}

	updatedRemainder, err := cfg.dbQ.UpdateRemainder(r.Context(), database.UpdateRemainderParams{
		Subject:     reqObj.Subject,
		Description: reqObj.Description,
		Timing:      theTime,
		HasPriority: reqObj.Has_Priority,
		Userid:      user.ID,
		ID:          id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			replyWithError("No remainder found", 200, w)
		} else {
			replyWithError(fmt.Sprintf("err in UpdateRemainder -> %v", err), 500, w)
		}
		return
	}

	type respStruct struct {
		Success   bool      `json:"success"`
		Remainder Remainder `json:"remainder"`
	}

	replyWithJson(w, respStruct{
		Success:   true,
		Remainder: dbToRespRemainder(updatedRemainder),
	}, 200)
}

func (cfg *apiConfig) deleteRemainder(w http.ResponseWriter, r *http.Request, user database.User) {
	// get url params
	idStr := chi.URLParam(r, "ID")

	id, err := uuid.Parse(idStr)
	if err != nil {
		replyWithError("not a valid remainder id", 400, w)
		return
	}

	deletedRemainder, err := cfg.dbQ.DeleteRemainder(r.Context(), database.DeleteRemainderParams{
		ID:     id,
		Userid: user.ID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			replyWithError("No remainder found", 200, w)
		} else {
			replyWithError(fmt.Sprintf("err in DeleteRemainder -> %v", err), 500, w)
		}
		return
	}

	type respStruct struct {
		Success   bool      `json:"success"`
		Remainder Remainder `json:"deleted_remainder"`
	}

	replyWithJson(w, respStruct{
		Success:   true,
		Remainder: dbToRespRemainder(deletedRemainder),
	}, 200)
}
