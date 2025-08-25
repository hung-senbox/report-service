package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReportHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ReportID  primitive.ObjectID `bson:"report_id"`
	EditorID  string             `bson:"editor_id"`
	Report    *Report            `bson:"report"`
	Timestamp int64              `bson:"timestamp"`
}
