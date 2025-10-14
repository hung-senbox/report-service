package request

type GetClassroomReportRequest4Web struct {
	TeacherID     string `json:"teacher_id" binding:"required"`
	TopicID       string `json:"topic_id" binding:"required"`
	TermID        string `json:"term_id" binding:"required"`
	UniqueLangKey string `json:"unique_lang_key" binding:"required"`
	ClassroomID   string `json:"classroom_id" binding:"required"`
}
