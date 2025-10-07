package request

type GetReportRequest4App struct {
	StudentID string `json:"student_id" binding:"required"`
	TopicID   string `json:"topic_id" binding:"required"`
	TermID    string `json:"term_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
}

type GetReportRequest4Web struct {
	StudentID string `json:"student_id" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
	TopicID   string `json:"topic_id" binding:"required"`
	TermID    string `json:"term_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
}
