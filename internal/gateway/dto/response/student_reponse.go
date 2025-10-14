package dto

type StudentResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
}

type Student4ClassroomReport struct {
	StudentID      string `json:"user_id"`
	StudentName    string `json:"user_name"`
	OrganizationID string `json:"organization_id"`
}
