package dto

type GetClassroomAssignTemplate struct {
	ClassroomID     string           `json:"class_id"`
	ClassroomName   string           `json:"class_name"`
	ClassroomIcon   string           `json:"class_icon"`
	Leader          Leader           `json:"leader"`
	AssignTemplates []AssignTemplate `json:"assign_templates"`
}

type AssignTemplate struct {
	TeacherID string `json:"teacher_id"`
	StudentID string `json:"student_id"`
}

type Leader struct {
	LeaderID   string `json:"leader_id"`
	LeaderRole string `json:"owner_role"`
	Avatar     Avatar `json:"avatar"`
}
