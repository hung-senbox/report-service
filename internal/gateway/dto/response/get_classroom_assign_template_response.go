package dto

type GetClassroomAssignTemplate struct {
	ClassroomID     string           `json:"classroom_id"`
	ClassroomName   string           `json:"classroom_name"`
	ClassroomIcon   string           `json:"class_icon"`
	AssignTemplates []AssignTemplate `json:"assign_templates"`
}

type AssignTemplate struct {
	TeacherID string `json:"teacher_id"`
	StudentID string `json:"student_id"`
}
