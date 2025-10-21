package response

import (
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/internal/report/model"
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
	Editing                    bool                        `json:"editing"`
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
	Now                 string `json:"now"`
	NowUpdatedAt        string `json:"now_updated_at"`
	Conclusion          string `json:"conclusion"`
	ConclusionUpdatedAt string `json:"conclusion_updated_at"`
	TermTitle           string `json:"term_title"`
}

type TeacherReportPreviousTerm struct {
	Now                 string `json:"now"`
	NowUpdatedAt        string `json:"now_updated_at"`
	Conclusion          string `json:"conclusion"`
	ConclusionUpdatedAt string `json:"conclusion_updated_at"`
	TermTitle           string `json:"term_title"`
}

type ClassroomReportResponse4Web struct {
	Student StudentReportClassroom `json:"student"`
	Teacher TeacherReportClassroom `json:"teacher"`
	Report  ReportResponse         `json:"report"`
}

type GetClassroomReportResponse4Web struct {
	Reports          []ClassroomReportResponse4Web `json:"reports"`
	SchoolTemplate   model.Template                `json:"school_template"`
	ClassroomTempate model.Template                `json:"classroom_template"`
}

type StudentReportClassroom struct {
	StudentID      string `json:"id"`
	StudentName    string `json:"name"`
	OrganizationID string `json:"organization_id"`
	AvatarMainUrl  string `json:"avatar_main_url"`
}

type TeacherReportClassroom struct {
	TeacherID      string `json:"id"`
	TeacherName    string `json:"name"`
	OrganizationID string `json:"organization_id"`
	AvatarMainUrl  string `json:"avatar_main_url"`
}
