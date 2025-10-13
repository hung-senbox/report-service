package model

type ReportPlanTemplate struct {
	ID             string   `json:"id" bson:"_id"`
	Template       Template `json:"template" bson:"template"`
	OrganizationID string   `json:"organization_id" bson:"organization_id"`
	TopicID        string   `json:"topic_id" bson:"topic_id"`
	TermID         string   `json:"term_id" bson:"term_id"`
	Language       string   `json:"language" bson:"language"`
	IsSchool       bool     `json:"is_school" bson:"is_school"`
	CreatedAt      int64    `json:"created_at" bson:"created_at"`
	UpdatedAt      int64    `json:"updated_at" bson:"updated_at"`
}

type Template struct {
	Title          string `json:"title" bson:"title"`
	CurriculumArea string `json:"curriculum_area" bson:"curriculum_area"`
	Introduction   string `json:"introduction" bson:"introduction"`
}
