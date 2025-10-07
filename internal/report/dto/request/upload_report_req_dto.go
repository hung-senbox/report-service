package request

type UploadReport4AppRequest struct {
	StudentID  string                 `json:"student_id" binding:"required"`
	TopicID    string                 `json:"topic_id" binding:"required"`
	TermID     string                 `json:"term_id" binding:"required"`
	Language   string                 `json:"language" binding:"required"`
	Status     string                 `json:"status" binding:"required"`
	Note       map[string]interface{} `json:"note" binding:"required"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}

type UploadReport4AWebRequest struct {
	StudentID  string                 `json:"student_id" binding:"required"`
	TopicID    string                 `json:"topic_id" binding:"required"`
	TermID     string                 `json:"term_id" binding:"required"`
	Language   string                 `json:"language" binding:"required"`
	Status     string                 `json:"status" binding:"required"`
	Note       map[string]interface{} `json:"note" binding:"required"`
	ReportData map[string]interface{} `json:"report_data" binding:"required"`
}

type NoteRequest struct {
	Color          string `json:"color"`
	ManagerComment string `json:"manager_comment"`
	ManagerNote    string `json:"manager_note"`
	NoteForTeacher string `json:"note_for_teacher"`
	TeacherReport  string `json:"teacher_report"`
	Status         string `json:"status"`
}

type ReportDataRequest struct {
	Before       Before       `json:"before"`
	Now          Now          `json:"now"`
	Conclusion   Conclusion   `json:"conclusion"`
	Goal         Goal         `json:"goal"`
	Introduction Introduction `json:"introduction"`
	PreviousTerm PreviousTerm `json:"previous_term"`
	Title        Title        `json:"title"`
	SubTitle     SubTitle     `json:"sub_title"`
}

type Before struct {
	Color          string `json:"color"`
	ManagerNote    string `json:"manager_note"`
	ManagerComment string `json:"manager_comment"`
	NoteForTeacher string `json:"note_for_teacher"`
	TeacherReport  string `json:"teacher_report"`
	Status         string `json:"status"`
}

type Now struct {
	Color          string `json:"color"`
	ManagerNote    string `json:"manager_note"`
	ManagerComment string `json:"manager_comment"`
	NoteForTeacher string `json:"note_for_teacher"`
	TeacherReport  string `json:"teacher_report"`
	Status         string `json:"status"`
}

type Conclusion struct {
	Color          string `json:"color"`
	ManagerNote    string `json:"manager_note"`
	ManagerComment string `json:"manager_comment"`
	NoteForTeacher string `json:"note_for_teacher"`
	TeacherReport  string `json:"teacher_report"`
	Status         string `json:"status"`
}

type Goal struct {
	Content string `json:"content"`
}

type Introduction struct {
	Color          string `json:"color"`
	ManagerNote    string `json:"manager_note"`
	ManagerComment string `json:"manager_comment"`
	NoteForTeacher string `json:"note_for_teacher"`
	TeacherReport  string `json:"teacher_report"`
	Status         string `json:"status"`
}

type PreviousTerm struct {
	Content string `json:"content"`
}

type SubTitle struct {
	Content string `json:"content"`
}

type Title struct {
	Content string `json:"content"`
}
