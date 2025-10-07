package mapper

import (
	gw_response "report-service/internal/gateway/dto/response"
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

func MapReportToResDTO(report *model.Report, teacher *gw_response.TeacherResponse) response.ReportResponse {

	return response.ReportResponse{
		ID:         report.ID.Hex(),
		StudentID:  report.StudentID,
		TopicID:    report.TopicID,
		TermID:     report.TermID,
		Language:   report.Language,
		Status:     report.Status,
		Note:       report.Note,
		ReportData: report.ReportData,
		CreatedAt:  report.CreatedAt,
		Editor:     *teacher,
	}
}

// MapReportListToResDTO maps slice of model.Report to slice of ReportResponse
func MapReportListToResDTO(reports []*model.Report) []response.ReportResponse {
	result := make([]response.ReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, MapReportToResDTO(r, nil))
	}
	return result
}
