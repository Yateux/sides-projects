package models

import (
	"time"
	_ "gopkg.in/mgo.v2/bson"
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
