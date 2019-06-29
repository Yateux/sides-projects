package utils

import (
	"encoding/json"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var DbName = os.Getenv("MONGO_DB_NAME")

// generic function to respond to a request
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetMSDatabase(dbSession *mgo.Session) *mgo.Database {
	return dbSession.DB("TOSET")
}

// empty response
func RespondEmpty(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write(nil)
}

// specific function to respond with an error
func RespondWithError(w http.ResponseWriter, code int, errorMessage string) {
	RespondWithJSON(w, code, map[string]interface{}{"error": errorMessage})
}

// check the method in the request to see if it is part of the allowed method for a route
func CheckHttpMethod(r *http.Request, allowedMethods []string) bool {
	for _, allowedMethod := range allowedMethods {
		if allowedMethod == r.Method {
			return true
		}
	}
	return false
}
