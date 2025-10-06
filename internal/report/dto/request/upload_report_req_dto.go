package request

type UploadReport4AppRequest struct {
	StudentID  string                 `json:"student_id" binding:"required"`
	TopicID    string                 `json:"topic_id" binding:"required"`
	TermID     string                 `json:"term_id" binding:"required"`
	Language   string                 `json:"language" binding:"required"`
	Status     string                 `json:"status" binding:"required"`
	Note       map[string]interface{} `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}

type UploadReport4AWebRequest struct {
	StudentID  string                 `json:"student_id" binding:"required"`
	TopicID    string                 `json:"topic_id" binding:"required"`
	TermID     string                 `json:"term_id" binding:"required"`
	Language   string                 `json:"language" binding:"required"`
	Status     string                 `json:"status" binding:"required"`
	Note       map[string]interface{} `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}
