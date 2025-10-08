package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReportHistory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ReportID   primitive.ObjectID `bson:"report_id"`
	EditorID   string             `bson:"editor_id"`
	EditorRole string             `bson:"editor_role"`
	Report     *Report            `bson:"report"`
	Timestamp  int64              `bson:"timestamp"`
}
