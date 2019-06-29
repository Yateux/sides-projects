package _auth

import (
	fabric "collections/fabric"
	colType "collections/types"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	test "./tests"
	"gopkg.in/mgo.v2/bson"
)

type responseJson struct {
	Error  string
	UserId string
	Token  string
}

func readJsonError(r io.Reader) string {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err.Error())
	}
	var e responseJson
	err = json.Unmarshal(body, &e)
	if err != nil {
		fmt.Println("Error - readJsonError:", err)
	}
	return e.Error
}

func readResponseJson(r io.Reader) responseJson {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err.Error())
	}
	var output responseJson
	err = json.Unmarshal(body, &output)
	if err != nil {
		fmt.Println("Error - readJsonError:", err)
	}
	return output
}
func TestAuthNoBody(t *testing.T) {
	DbClean()
	r := test.CreateRequest("POST", "/v1/ops/authentication", nil, test.UserData{}, db)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Authentication(w, r)
	resp := w.Result()
	statusContent := readJsonError(w.Body)
	if resp.StatusCode != 500 || statusContent != "Decode body failed" {
		t.Error(fmt.Sprintf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m{\"error\":\"%s\"}\033[0m' not '\x1b[1;31m%s\033[0m'.", 500, resp.StatusCode, "Decode body failed", statusContent))
	}
}

func TestAuthWrongUser(t *testing.T) {
	DbClean()
	userTest := fabric.UserInit(db.C("users"), "admin")
	_, err := userTest.CChangeEmail("itIsMy@email.co").Save()
	if err != nil {
		t.Error(err)
	}
	body := []byte(`{"email": "wrongEmail", "token": "f2d81a260dea8a100dd517984e53c56a7523d96942a834b9cdc249bd4e8c7aa9"}`)
	r := test.CreateRequest("POST", "/v1/ops/authentication", body, test.UserData{}, db)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Authentication(w, r)
	resp := w.Result()
	statusContent := readJsonError(w.Body)
	if resp.StatusCode != 403 || statusContent != "User or password incorrect" {
		t.Error(fmt.Sprintf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", 403, resp.StatusCode, "User or password incorrect", statusContent))
	}
}

func TestAuthWrongPassword(t *testing.T) {
	DbClean()
	userTest := fabric.UserInit(db.C("users"), "admin")
	_, err := userTest.CChangeEmail("itIsMy@email.co").CChangePasswordAzerty().Save()
	if err != nil {
		t.Error(err)
	}
	body := []byte(`{"email": "itIsMy@email.co", "hashedPassword": "itsaWrongPassword"}`)
	r := test.CreateRequest("POST", "/v1/ops/authentication", body, test.UserData{}, db)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Authentication(w, r)
	resp := w.Result()
	statusContent := readJsonError(w.Body)
	if resp.StatusCode != 403 || statusContent != "User or password incorrect" {
		t.Error(fmt.Sprintf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m and status content '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", 403, resp.StatusCode, "User or password incorrect", statusContent))
	}
}

// Password:       azerty
// HashedPassword: f2d81a260dea8a100dd517984e53c56a7523d96942a834b9cdc249bd4e8c7aa9
// SHA256:         $2a$10$7U4ZCnFYqnKt8eS8B.i0ZO5tIua17CGn4084olicf4AaIOyE.Lm6m
func TestAuth(t *testing.T) {
	DbClean()
	userTest := fabric.UserInit(db.C("users"), "admin")
	userId, err := userTest.CChangeEmail("itIsMy@email.co").CChangePasswordAzerty().Save()
	if err != nil {
		t.Error(err)
	}
	body := []byte(`{"email": "itIsMy@email.co", "hashedPassword": "f2d81a260dea8a100dd517984e53c56a7523d96942a834b9cdc249bd4e8c7aa9"}`)
	r := test.CreateRequest("POST", "/v1/ops/authentication", body, test.UserData{}, db)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Authentication(w, r)
	resp := w.Result()
	statusContent := readResponseJson(w.Body)
	if resp.StatusCode != 200 || statusContent.UserId != userId || statusContent.Token == "" {
		t.Error(fmt.Sprintf("Must return an error with http code \x1b[1;32m%d\033[0m not \x1b[1;31m%d\033[0m, json userId '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m' and json userId '\x1b[1;32m%s\033[0m' not '\x1b[1;31m%s\033[0m'.", 200, resp.StatusCode, userId, statusContent.UserId, "[hash code]", statusContent.Token))
	}
	var user colType.User
	if err := db.C("users").Find(bson.M{
		"_id": userId,
	}).One(&user); err != nil {
		t.Error(err)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(statusContent.Token))
	hashStringNowInDatabase := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if user.Services.Resume.LoginTokens[0].When.IsZero() == true {
		t.Error("'When' set in the db must not be 'nil'")
	}
	if user.Services.Resume.LoginTokens[0].HashedToken != hashStringNowInDatabase {
		t.Error(fmt.Sprintf("'Token' set in the db is equal to '%s' but expect '%s'", user.Services.Resume.LoginTokens[0].HashedToken, hashStringNowInDatabase))
	}
	return
}
