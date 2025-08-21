package request

type UploadReportRequestDTO struct {
	Key        string                 `json:"key" binding:"required"`
	Note       string                 `json:"note,omitempty"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}
