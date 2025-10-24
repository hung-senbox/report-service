package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportTranslation struct {
	ID           primitive.ObjectID               `bson:"_id,omitempty"`
	StudentID    string                           `bson:"student_id"`
	TopicID      string                           `bson:"topic_id"`
	TermID       string                           `bson:"term_id"`
	Translations map[string]ReportTranslationData `bson:"translations"`
	CreatedAt    time.Time                        `bson:"created_at"`
	UpdatedAt    time.Time                        `bson:"updated_at"`
}

type ReportTranslationData struct {
	Before     string `json:"before"`
	Conclusion string `json:"conclusion"`
	Now        string `json:"now"`
}
