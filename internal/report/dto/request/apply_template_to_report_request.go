package request

type ApplyTemplateIsSchoolToReportRequest struct {
	TermID         string `json:"term_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	UniqueLangKey  string `json:"unique_lang_key" binding:"required"`
	Title          string `json:"title"`
	Introduction   string `json:"introduction"`
	CurriculumArea string `json:"curriculum_area"`
}

type ApplyTemplateIsClassroomToReportRequest struct {
	TermID         string `json:"term_id" binding:"required"`
	TopicID        string `json:"topic_id" binding:"required"`
	ClassroomID    string `json:"classroom_id" binding:"required"`
	UniqueLangKey  string `json:"unique_lang_key" binding:"required"`
	Title          string `json:"title"`
	Introduction   string `json:"introduction"`
	CurriculumArea string `json:"curriculum_area"`
}
