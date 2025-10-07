package mapper

import (
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

// MapReportHistoryToResDTO maps a single ReportHistory to ReportHistoryResponse
func MapReportHistoryToResDTO(history *model.ReportHistory) response.ReportHistoryResponse {
	var reportRes response.ReportResponse
	if history.Report != nil {
		reportRes = MapReportToResDTO(history.Report, nil)
	}

	return response.ReportHistoryResponse{
		ID:        history.ID.Hex(),
		ReportID:  history.ReportID.Hex(),
		EditorID:  history.EditorID,
		Report:    reportRes,
		Timestamp: history.Timestamp,
	}
}

// MapReportHistoryListToResDTO maps slice of ReportHistory to slice of ReportHistoryResponse
func MapReportHistoryListToResDTO(histories []*model.ReportHistory) []response.ReportHistoryResponse {
	result := make([]response.ReportHistoryResponse, 0, len(histories))
	for _, h := range histories {
		result = append(result, MapReportHistoryToResDTO(h))
	}
	return result
}
