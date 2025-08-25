package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StudentID  string             `bson:"student_id"`
	TopicID    string             `bson:"topic_id"`
	TermID     string             `bson:"term_id"`
	Language   string             `bson:"language"`
	Status     string             `bson:"status"`
	Note       bson.M             `bson:"note"`
	ReportData bson.M             `bson:"report_data" json:"report_data"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
