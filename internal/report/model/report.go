package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	EditorID   string             `bson:"editor_id,omitempty"`
	StudentID  string             `bson:"student_id"`
	TopicID    string             `bson:"topic_id"`
	TermID     string             `bson:"term_id"`
	Language   string             `bson:"language"`
	Status     string             `bson:"status"`
	Editing    *bool              `bson:"editing" json:"editing"`
	ReportData bson.M             `bson:"report_data" json:"report_data"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type ReportData struct {
	Before         Section `json:"before"`
	Conclusion     Section `json:"conclusion"`
	CurriculumArea Section `json:"curriculum_area"`
	Goal           Section `json:"goal"`
	Introduction   Section `json:"introduction"`
	Note           Section `json:"note"`
	Now            Section `json:"now"`
	PreviousTerm   Section `json:"previous_term"`
	SubTitle       Section `json:"sub_title"`
	Title          Section `json:"title"`
}

type Section struct {
	Color            string `json:"color"`
	Content          string `json:"content"`
	ManagerComment   string `json:"manager_comment"`
	ManagerNote      string `json:"manager_note"`
	ManagerUpdatedAt string `json:"manager_updated_at"`
	NoteForTeacher   string `json:"note_for_teacher"`
	Status           string `json:"status"`
	TeacherReport    string `json:"teacher_report"`
	UpdatedAt        string `json:"updated_at"`
}

type Reportstruct struct {
	ID         string     `json:"id"`
	StudentID  string     `json:"student_id"`
	TopicID    string     `json:"topic_id"`
	TermID     string     `json:"term_id"`
	EditorID   string     `json:"editor_id"`
	Language   string     `json:"language"`
	Status     string     `json:"status"`
	ReportData ReportData `json:"report_data"`
	CreatedAt  time.Time  `json:"created_at"`
	Progress   int        `json:"progress"`
}
