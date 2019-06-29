package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	fileErr "../errors"

	raven "github.com/getsentry/raven-go"
	mgo "gopkg.in/mgo.v2"
)

type bodyContent struct {
	Email          string
	HashedPassword string
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_-"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func getRandomString(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func RandomId() string {
	bytes := make([]byte, 17)
	for i := 0; i < 17; i++ {
		r := rand.Intn(61)
		switch {
		case r < 10:
			r = r + 48
		case r < 36:
			r = r + 55
		default:
			r = r + 61
		}
		bytes[i] = byte(r)
	}
	return string(bytes)
}

func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check and return error in case of wrong method or DB error
func CheckMethodGetDB(method []string, r *http.Request) (*mgo.Database, int, string) {
	if ok := CheckHttpMethod(r, method); !ok {
		return nil, 404, "page not found"
	}
	db, ok := r.Context().Value("database").(*mgo.Database)
	if !ok {
		msg := fmt.Sprintf("[TOSET - DATABASE] %s %s - Database connection failed : %v", r.Method, r.URL.RequestURI(), ok)
		log.Println(msg)
		raven.CaptureError(fileErr.SideError(msg), nil)
		return nil, 500, "Problem with database connection"
	}
	return db, 0, ""
}

// Check and return error in case of wrong method, DB, dbName, dbSession error
func CheckMethodGetDBSession(method []string, r *http.Request) (*mgo.Database, *mgo.Session, int, string) {
	if ok := CheckHttpMethod(r, method); !ok {
		return nil, nil, 404, "page not found"
	}
	db, ok := r.Context().Value("database").(*mgo.Database)
	if !ok {
		msg := fmt.Sprintf("[TOSET - DATABASE] %s %s - Database connection failed : %v", r.Method, r.URL.RequestURI(), ok)
		log.Println(msg)
		raven.CaptureError(fileErr.SideError(msg), nil)
		return nil, nil, 500, "Problem with database connection"
	}
	dbSession, ok := r.Context().Value("dbSession").(*mgo.Session)
	if !ok {
		msg := fmt.Sprintf("[TOSET - DATABASE SESSION] %s %s - Database connection failed : %v", r.Method, r.URL.RequestURI(), ok)
		log.Println(msg)
		raven.CaptureError(fileErr.SideError(msg), nil)
		return nil, nil, 500, "Problem with database connection"
	}
	return db, dbSession, 0, ""
}

type HandleErrors struct {
	Error   bool
	Code    int
	Content string
}

func (e *HandleErrors) HTTPError(code int, content string) {
	e.Error = true
	e.Code = code
	e.Content = content
}

// Check and return error in case of wrong method or DB error
func (e *HandleErrors) HandleMethodGetDB(method []string, r *http.Request) *mgo.Database {
	if e.Error == true {
		return nil
	}
	if ok := CheckHttpMethod(r, method); !ok {
		e.HTTPError(404, "page not found")
		return nil
	}
	db, ok := r.Context().Value("database").(*mgo.Database)
	if !ok {
		msg := fmt.Sprintf("[TOSET - DATABASE] %s %s - Database connection failed : %v", r.Method, r.URL.RequestURI(), ok)
		log.Println(msg)
		raven.CaptureError(fileErr.SideError(msg), nil)
		e.HTTPError(500, "Problem with database connection")
		return nil
	}
	return db
}
