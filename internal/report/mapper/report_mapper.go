package mapper

import (
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

func MapReportToResDTO(report *model.Report, teacher *gw_response.TeacherResponse, managerCmPrevious response.ManagerCommentPreviousTerm, teacherRpPrevious response.TeacherReportPreviousTerm) response.ReportResponse {

	var teacherEditor gw_response.TeacherResponse
	if teacher != nil {
		teacherEditor = *teacher
	}
	return response.ReportResponse{
		ID:                         report.ID.Hex(),
		StudentID:                  report.StudentID,
		TopicID:                    report.TopicID,
		TermID:                     report.TermID,
		Language:                   report.Language,
		Status:                     report.Status,
		ReportData:                 report.ReportData,
		CreatedAt:                  report.CreatedAt,
		Editor:                     teacherEditor,
		ManagerCommentPreviousTerm: managerCmPrevious,
		TeacherReportPreviousTerm:  teacherRpPrevious,
	}
}

// MapReportListToResDTO maps slice of model.Report to slice of ReportResponse
func MapReportListToResDTO(reports []*model.Report) []response.ReportResponse {
	result := make([]response.ReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, MapReportToResDTO(r, nil, response.ManagerCommentPreviousTerm{}, response.TeacherReportPreviousTerm{}))
	}
	return result
}
