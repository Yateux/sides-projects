package _auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"../utils"

	colType "collections/types"

	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func generateRandomSHA256() (string, string) {
	hash := sha256.New()
	generated := getRandomString(43)
	hash.Write([]byte(generated))
	hashString := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hashString, generated
}

func Authentication(w http.ResponseWriter, r *http.Request) {
	// Check and return error in case of wrong method or db error connection
	db, errCode, errContent := CheckMethodGetDB([]string{"POST"}, r)
	if errCode != 0 && errContent != "" {
		utils.RespondWithError(w, errCode, errContent)
		return
	}
	// Manage body
	decoder := json.NewDecoder(r.Body)
	var body bodyContent
	if err := decoder.Decode(&body); err != nil {
		log.Println("[DECODE BODY] /v1/authentication :", err)
		utils.RespondWithError(w, 500, "Decode body failed")
		return
	}
	// Find the user in the db with email from the body
	var user colType.User
	if err := db.C("users").Find(bson.M{
		"emails.address": body.Email,
	}).One(&user); err != nil {
		log.Println("[DB FIND] /v1/authentication :", err)
		if err == mgo.ErrNotFound {
			utils.RespondWithError(w, 403, "User or password incorrect")
			return
		}
		utils.RespondWithError(w, 500, "Error connecting with database")
		return
	}
	// Comparing the password with the hashed password from the body
	err := bcrypt.CompareHashAndPassword([]byte(user.Services.Password.BCrypt), []byte(body.HashedPassword))
	if err != nil {
		log.Println("[HASH ERROR] /v1/authentication :", err)
		utils.RespondWithError(w, 403, "User or password incorrect")
		return
	}
	today := time.Now()
	// Create a 43 characters random string then hash with SHA256
	hashToken, token := generateRandomSHA256()
	// Update the db with the new token
	newElementToken := colType.LoginToken{
		When:        today,
		HashedToken: hashToken,
	}
	if err := db.C("users").Update(bson.M{
		"_id": user.ID,
	}, bson.M{
		"$push": bson.M{
			"services.resume.loginTokens": newElementToken,
		},
	}); err != nil {
		log.Println("[DB UPDATE] /v1/authentication :", err)
		utils.RespondWithError(w, 500, "Error updating the database")
		return
	}
	utils.RespondWithJSON(w, 200, map[string]interface{}{
		"userId": user.ID,
		"token":  token,
	})
}
