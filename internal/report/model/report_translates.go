package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportTranslation struct {
	ID           primitive.ObjectID               `json:"id" bson:"_id,omitempty"`
	StudentID    string                           `json:"student_id" bson:"student_id"`
	TopicID      string                           `json:"topic_id" bson:"topic_id"`
	TermID       string                           `json:"term_id" bson:"term_id"`
	Translations map[string]ReportTranslationData `json:"translations" bson:"translations"`
	CreatedAt    time.Time                        `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time                        `json:"updated_at" bson:"updated_at"`
}

type ReportTranslationData struct {
	Before     string `json:"before"`
	Conclusion string `json:"conclusion"`
	Now        string `json:"now"`
}
