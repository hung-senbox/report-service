package response

import "report-service/pkg/constants"

type GetTeacherReportTasksResponse4App struct {
	Term        string                      `json:"term"`
	Topic       string                      `json:"topic"`
	StudentName string                      `json:"student_name"`
	Deadline    string                      `json:"deadline"`
	Task        constants.TeacherReportTask `json:"task"`
	Status      string                      `json:"status"`
	Language    string                      `json:"language"`
}
