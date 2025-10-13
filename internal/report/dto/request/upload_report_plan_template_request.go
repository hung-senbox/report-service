package request

type UploadReportPlanTemplateRequest struct {
	TopicID        string `json:"topic_id" binding:"required"`
	TermID         string `json:"term_id" binding:"required"`
	Language       string `json:"language" binding:"required"`
	Goal           string `json:"goal" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Introduction   string `json:"introduction" binding:"required"`
	CurriculumArea string `json:"curriculum_area" binding:"required"`
	IsSchool       bool   `json:"is_school"`
}
