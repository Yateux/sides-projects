package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	mgo "gopkg.in/mgo.v2"
)

func CreateRequest(method string, url string, body []byte, userData UserData, db *mgo.Database) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), "userId", userData.UserId)
	ctx = context.WithValue(ctx, "hashedToken", userData.HashedToken)
	ctx = context.WithValue(ctx, "database", db)
	return r.WithContext(ctx)
}

func CreateRequestSession(method string, url string, body []byte, userData UserData, db *mgo.Database, dbSession *mgo.Session) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), "userId", userData.UserId)
	ctx = context.WithValue(ctx, "hashedToken", userData.HashedToken)
	ctx = context.WithValue(ctx, "database", db)
	ctx = context.WithValue(ctx, "dbSession", dbSession)
	return r.WithContext(ctx)
}

func CreateRequestWithoutDB(method string, url string, body []byte, userData UserData) *http.Request {
	r := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	ctx := context.WithValue(r.Context(), "userId", userData.UserId)
	ctx = context.WithValue(ctx, "hashedToken", userData.HashedToken)
	return r.WithContext(ctx)
}

func ChargeResponse(w *httptest.ResponseRecorder, response interface{}) error {
	res := w.Result()
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(response)
	return err
}
