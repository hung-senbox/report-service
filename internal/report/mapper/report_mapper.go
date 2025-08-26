package mapper

import (
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

func MapReportToResDTO(report *model.Report) response.ReportResponse {
	// get editor info
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
	}
}

// MapReportListToResDTO maps slice of model.Report to slice of ReportResponse
func MapReportListToResDTO(reports []*model.Report) []response.ReportResponse {
	result := make([]response.ReportResponse, 0, len(reports))
	for _, r := range reports {
		result = append(result, MapReportToResDTO(r))
	}
	return result
}
