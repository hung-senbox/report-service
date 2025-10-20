package request

type GetClassroomReportRequest4Web struct {
	TopicID       string `json:"topic_id" binding:"required"`
	TermID        string `json:"term_id" binding:"required"`
	UniqueLangKey string `json:"unique_lang_key" binding:"required"`
	ClassroomID   string `json:"classroom_id" binding:"required"`
}
