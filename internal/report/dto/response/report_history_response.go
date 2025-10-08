package response

import "time"

type ReportHistoryResponse struct {
	ID        string         `json:"id"`
	ReportID  string         `json:"report_id"`
	EditorID  string         `json:"editor_id"`
	Report    ReportResponse `json:"report,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}
