package response

import "time"

type ReportResponse struct {
	ID         string                 `json:"id"`
	StudentID  string                 `json:"student_id"`
	TopicID    string                 `json:"topic_id"`
	TermID     string                 `json:"term_id"`
	Language   string                 `json:"language"`
	Status     string                 `json:"status"`
	Note       map[string]interface{} `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}
