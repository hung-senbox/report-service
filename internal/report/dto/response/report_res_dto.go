package response

import (
	"report-service/internal/gateway/dto"
	"time"
)

type ReportResponse struct {
	ID         string                 `json:"id"`
	StudentID  string                 `json:"student_id"`
	TopicID    string                 `json:"topic_id"`
	TermID     string                 `json:"term_id"`
	Editor     dto.TeacherResponse    `json:"editor"`
	Language   string                 `json:"language"`
	Status     string                 `json:"status"`
	Note       map[string]interface{} `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data"`
	CreatedAt  time.Time              `json:"created_at"`
}

type ReportEditor struct {
	ID     string `json:"id"`
	Name   string `json:"full_name"`
	Avatar string `json:"avatar"`
}
