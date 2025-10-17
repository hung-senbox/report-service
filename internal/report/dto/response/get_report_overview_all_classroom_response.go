package response

type AllClassroomTopicStatus struct {
	TopicID    string `json:"topic_id"`
	Progress   int    `json:"progress"`
	Before     int    `json:"before"`
	Now        int    `json:"now"`
	Conclusion int    `json:"conclusion"`
}

type AllClassroomReport struct {
	ClassName string                    `json:"class_name"`
	DOB       string                    `json:"dob"`
	Age       int                       `json:"age"`
	Class     float64                   `json:"class"`
	Topics    []AllClassroomTopicStatus `json:"topics"` // ← đổi sang slice
}

type GetReportOverviewAllClassroomResponse struct {
	Reports []AllClassroomReport `json:"reports"`
}
