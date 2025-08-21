package response

import "time"

type ReportResponse struct {
	ID         string                 `json:"id"`
	Key        string                 `json:"key"`
	Note       string                 `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}
