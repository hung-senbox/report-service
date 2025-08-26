package dto

type TeacherResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Avatar         Avatar `json:"avatar"`
}
