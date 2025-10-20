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
	Image          string `json:"image"`
}

type AssingmentTemplate struct {
	Students []StudentTemplate `json:"students"`
	Teachers []TeacherTemplate `json:"teachers"`
}

type StudentTemplate struct {
	StudentID      string `json:"user_id"`
	StudentName    string `json:"user_name"`
	OrganizationID string `json:"organization_id"`
	Avatar         Avatar `json:"avatar"`
}

type TeacherTemplate struct {
	TeacherID      string `json:"user_id"`
	TeacherName    string `json:"user_name"`
	OrganizationID string `json:"organization_id"`
	Avatar         Avatar `json:"avatar"`
}
