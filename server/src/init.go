package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"./utils"
	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	mgo "gopkg.in/mgo.v2"
)

type Adapter func(http.Handler) http.Handler

var clusterEnv = os.Getenv("CLUSTER")

// initSentry user the CLUSTER environment variable to set the correct sentry DNS.
func initSentry() error {
	var sentryDns string
	if clusterEnv == "prod" {
		sentryDns = os.Getenv("SENTRY_GO_TOSET")
	} else {
		sentryDns = os.Getenv("SENTRY_TEST")
	}
	if err := raven.SetDSN(sentryDns); err != nil {
		return err
	}
	return nil
}

// dbInit launch the connection to the database. For preprod and prod environment, authenticate.
func dbInit() *mgo.Session {
	db_url := os.Getenv("MONGO_URL")
	log.Println("----- START CONFIGURATION DB ---")
	log.Println(os.Getenv("MONGO_URL"))
	log.Println(os.Getenv("MONGO_USER"))
	log.Println(os.Getenv("MONGO_PASSWORD"))
	log.Println(os.Getenv("MONGO_DB_NAME"))
	log.Println("----- END CONFIGURATION DB -----")
	fmt.Printf("mongo url: %s\n", db_url)
	db, err := mgo.Dial(db_url)
	
	if err != nil {

		raven.CaptureErrorAndWait(err, nil)
		log.Fatal("cannot connect to mongo\n", err)
	}
	
	// Careful here, the user and and the password depend on the database you use
	// the following vars strike on our main db
	if db_password := os.Getenv("MONGO_PASSWORD"); db_password != "" {
		if db_user := os.Getenv("MONGO_USER"); db_user != "" {
			if err := db.DB(utils.DbName).Login(db_user, db_password); err != nil {
				raven.CaptureErrorAndWait(err, nil)
				log.Fatal(err)
			}
			
		}
		
	}
	
	return db
}

// adapt transforms an handler without changing it's type. Usefull for authentification.
func adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// adapt the request by checking the auth and filling the context with usefull data
func enhanceHandlers(r *mux.Router, db *mgo.Session) http.Handler {
	return adapt(r, withRights(), withDB(db), withCors())
}

// withDB is an adapter that copy the access to the database to serve a specific call
// don't forget to change the db name here
func withDB(db *mgo.Session) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dbsession := db.Copy()
			defer dbsession.Close() // cleaning up
			db_name := os.Getenv("MONGO_DB_NAME")
			if db_name == "" {
				db_name = "TOSET"
			}
			ctx := context.WithValue(r.Context(), "database", dbsession.DB(db_name))
			ctx = context.WithValue(ctx, "dbSession", dbsession)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// here you have to chose one of the functions withRights. One works withDB
// an auth route, the other one works an apiKey

// don't forget you have to declare your apiKey
func withRights() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header
			apiKey := header.Get("apikey")
			if apiKey != os.Getenv("YOUR_SERVICE_API_KEY") {
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				fmt.Println("unauthorized")
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// withRights is an adapter that verify the user exists, verify the token, and attach userId and
// activeOrganisationId as organisationId to the request.
// works with an authorised route which is not automatic in the case of a micro service
// func withRights() Adapter {
// 	return func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			routeURL := *r.URL
// 			if routeURL.String() == "/v1/authentication" {
// 				h.ServeHTTP(w, r)
// 				return
// 			}
// 			// retrieve the data send on every request. adminId serves for logas purpose. (not implemented yet)
// 			// vars := mux.Vars(r)
// 			userId := r.Header.Get("userId")
// 			token := r.Header.Get("token")
// 			// hash the token to compare it to the one on the db.
// 			hashing := sha256.New()
// 			hashing.Write([]byte(token))
// 			hashedToken := base64.StdEncoding.EncodeToString(hashing.Sum(nil))
// 			db, ok := r.Context().Value("database").(*mgo.Database)
// 			if !ok {
// 				utils.RespondWithError(w, 500, "Problem with database connection")
// 				return
// 			}
// 			var err error
// 			// Find the admin user
// 			var user colType.User
// 			err = db.C("users").Find(bson.M{
// 				"_id": userId,
// 				"roles.__global_roles__": bson.M{
// 					"$in": []interface{}{
// 						"maestro",
// 						"admin",
// 					},
// 				},
// 			}).One(&user)
// 			if err != nil {
// 				if err == mgo.ErrNotFound {
// 					fmt.Println("[DB REQUEST] Error: This user doesn't exists or is not admin")
// 					utils.RespondWithError(w, 403, "Access forbidden")
// 					return
// 				}
// 				fmt.Println("[DB REQUEST] Error: Mongodb request failed")
// 				utils.RespondWithError(w, 500, "Error in fetching user: "+err.Error())
// 				return
// 			}
// 			// Exception if routeURL is authentication because token doesn't exists
// 			// check that this user has this token
// 			hasToken := false
// 			for _, loginToken := range user.Services.Resume.LoginTokens {
// 				if hashedToken == loginToken.HashedToken {
// 					hasToken = true
// 				}
// 			}
// 			if !hasToken {
// 				fmt.Println("token does not exist")
// 				utils.RespondWithError(w, 403, "forbidden")
// 				return
// 			}
// 			// attach usefull data to the request
// 			ctx := context.WithValue(r.Context(), "userId", user.ID)
// 			h.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// withCors is an adpater that allowed the specific headers we need for our requests from a
// different domain.
func withCors() Adapter {
	return func(h http.Handler) http.Handler {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost"},
			AllowedHeaders:   []string{"userId", "token", "Content-type"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowCredentials: true,
		})
		return c.Handler(h)
	}
}
