package request

type ApplyTemplateIsSchoolToReportRequest struct {
	TermID         string `json:"term_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	UniqueLangKey  string `json:"unique_lang_key" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Introduction   string `json:"introduction" binding:"required"`
	CurriculumArea string `json:"curriculum_area" binding:"required"`
	IsChool        *bool  `json:"is_school" binding:"required"`
}

type ApplyTemplateIsClassroomToReportRequest struct {
	TermID         string `json:"term_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	ClassroomID    string `json:"classroom_id" binding:"required"`
	UniqueLangKey  string `json:"unique_lang_key" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Introduction   string `json:"introduction" binding:"required"`
	CurriculumArea string `json:"curriculum_area" binding:"required"`
	IsChool        *bool  `json:"is_school" binding:"required"`
}
