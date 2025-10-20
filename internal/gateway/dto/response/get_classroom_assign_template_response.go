package dto

type GetAllClassroomAssignTemplate struct {
	ClassroomID      string           `json:"classroom_id"`
	ClassroomName    string           `json:"classroom_name"`
	ClassroomIconUrl string           `json:"classroom_icon_url"`
	AssignTemplates  []AssignTemplate `json:"assign_templates"`
}

type AssignTemplate struct {
	TeacherID string `json:"teacher_id"`
	StudentID string `json:"student_id"`
}
