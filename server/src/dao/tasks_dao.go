package dao

import (


	"../utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TasksDAO struct {
	Server   string
	Database string
}


const (
	COLLECTION = "tasks"
)

// Find list of movies
func (m *TasksDAO) FindAll() ([]Tasks, error) {
	var tasks []Tasks
	err := db.C(COLLECTION).Find(bson.M{}).All(&tasks)
	return tasks, err
}

// Update an existing movie
func (m *TasksDAO) Update(tasks Tasks) error {
	err := db.C(COLLECTION).UpdateId(tasks.ID, &tasks)
	return err
}