package mapper

import (
	"report-service/internal/report/dto/response"
	"report-service/internal/report/model"
)

// MapReportHistoryToResDTO maps a single ReportHistory to ReportHistoryResponse
func MapReportHistoryToRes4App(history *model.ReportHistory) response.ReportHistoryResponse {
	var reportRes response.ReportResponse
	if history.Report != nil {
		reportRes = MapReportToResDTO(history.Report, nil, response.ManagerCommentPreviousTerm{}, response.TeacherReportPreviousTerm{}, "")
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
func MapReportHistoryListToRes4App(histories []*model.ReportHistory) []response.ReportHistoryResponse {
	result := make([]response.ReportHistoryResponse, 0, len(histories))
	for _, h := range histories {
		result = append(result, MapReportHistoryToRes4App(h))
	}
	return result
}
