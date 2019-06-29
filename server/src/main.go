package main

import (
	"encoding/json"
	"flag"

	"fmt"
	"log"
	"net/http"
	"time"

	"./utils"

	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Tasks struct {
	ID        string                 `bson:"_id" json:"id"`
	Name      string                 `bson:"name" json:"name"`
	Visible   bool                   `bson:"visible" json:"visible"`
	Country   string                 `bson:"country" json:"country"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
	Type      map[string]interface{} `bson:”type” json:”type”`
	Shifts    []interface{}          `bson:”shifts” json:”shifts”`
}

type Freelancers struct {
	ID   string        `bson:"_id" json:"id"`
	Tags []interface{} `bson:”tags” json:”tags”`
}

// GET list of movies
func FindTasksEndpoint(w http.ResponseWriter, r *http.Request) {
	var tasks []Tasks
	dbs, _, _ := utils.CheckMethodGetDB([]string{"GET"}, r)
	err := dbs.C("tasks").Find(bson.M{}).All(&tasks)

	if err != nil {
		log.Println("error")
		return
	}

	utils.RespondWithJSON(w, 200, tasks)
}



func UpdateTasksEndPoint(w http.ResponseWriter, r *http.Request) {

	var tasks Tasks
	dbs, _, _ := utils.CheckMethodGetDB([]string{"PUT"}, r)

	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := dbs.C("tasks").Update(bson.M{"_id": tasks.ID}, bson.M{"$set": bson.M{"visible": tasks.Visible}}); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success", "tasks": tasks.ID})
}


// GET list of freelancers
func FindFreelancersEndpoint(w http.ResponseWriter, r *http.Request) {
	var freelancers []Freelancers
	dbs, _, _ := utils.CheckMethodGetDB([]string{"GET"}, r)
	err := dbs.C("freelancers").Find(bson.M{}).All(&freelancers)

	if err != nil {
		log.Println("error")
		return
	}

	utils.RespondWithJSON(w, 200, freelancers)
}

// Function that instantiates and populates the router
func Handlers() *mux.Router {
	// instantiating the router
	api := mux.NewRouter()

	api.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, "OK")
	})

	api.HandleFunc("/api/v1/tasks", FindTasksEndpoint).Methods("GET")
	api.HandleFunc("/api/v1/tasks", UpdateTasksEndPoint).Methods("PUT")
	api.HandleFunc("/api/v1/freelancers", FindFreelancersEndpoint).Methods("GET")

	// By default this route is not used. If you want to use, decomment withRights in init.go
	// also, auth.go is in folder _auth which is ignore by the build, if you want to use it, remove the _
	//api.HandleFunc("/v1/authentication", routes.Authentication).Methods("POST") // Don't forget the exception in init.go
	return api
}

func main() {
	// parsing flags
	portPtr := flag.String("port", "8080", "port your want to listen on")
	sentry := flag.Bool("sentry-mode", true, "should sentry be loaded and tracking")
	fmt.Println("environment: ", clusterEnv)
	flag.Parse()

	// sentry initialisation
	if *sentry {
		if err := initSentry(); err != nil {
			log.Fatal(err)
		}
	}

	if *portPtr != "" {

		fmt.Printf("running on port: %s\n", *portPtr)
	}

	db := dbInit()
	r := Handlers()
	enhancedRouter := enhanceHandlers(r, db)
	if err := http.ListenAndServe(":"+*portPtr, enhancedRouter); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}
}
