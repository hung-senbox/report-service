package mapper

import (
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

func MapReportToResDTO(report *model.Report) response.ReportResponse {
	return response.ReportResponse{
		ID:         report.ID.Hex(),
		Key:        report.Key,
		Note:       report.Note,
		ReportData: report.ReportData, // gán trực tiếp, vì đã là bson.M = map
		CreatedAt:  report.CreatedAt,
		UpdatedAt:  report.UpdatedAt,
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
