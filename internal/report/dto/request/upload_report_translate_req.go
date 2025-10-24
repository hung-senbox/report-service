package request

type UploadReportTranslateRequest struct {
	StudentID string                 `json:"student_id" binding:"required"`
	TopicID   string                 `json:"topic_id" binding:"required"`
	TermID    string                 `json:"term_id" binding:"required"`
	Language  string                 `json:"language" binding:"required"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}