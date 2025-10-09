package response

import (
	gw_response "report-service/internal/gateway/dto/response"
	"time"
)

type ReportResponse struct {
	ID                         string                      `json:"id"`
	StudentID                  string                      `json:"student_id"`
	TopicID                    string                      `json:"topic_id"`
	TermID                     string                      `json:"term_id"`
	Editor                     gw_response.TeacherResponse `json:"editor,omitempty"`
	Language                   string                      `json:"language"`
	Status                     string                      `json:"status"`
	ReportData                 map[string]interface{}      `json:"report_data"`
	CreatedAt                  time.Time                   `json:"created_at"`
	ManagerCommentPreviousTerm ManagerCommentPreviousTerm  `json:"manager_comment_previous_term"`
	TeacherReportPreviousTerm  TeacherReportPreviousTerm   `json:"teacher_report_previous_term"`
	LatestDataTermID           string                      `json:"latest_data_term_id"`
}

type ReportEditor struct {
	ID     string `json:"id"`
	Name   string `json:"full_name"`
	Avatar string `json:"avatar"`
}

type ManagerCommentPreviousTerm struct {
	Now        string `json:"now"`
	Conclusion string `json:"conclusion"`
	TermTitle  string `json:"term_title"`
}

type TeacherReportPreviousTerm struct {
	Now        string `json:"now"`
	Conclusion string `json:"conclusion"`
	TermTitle  string `json:"term_title"`
}
