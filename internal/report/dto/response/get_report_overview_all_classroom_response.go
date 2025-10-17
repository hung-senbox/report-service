package response

type AllClassroomTopicStatus struct {
	Progress   int `json:"progress"`
	Before     int `json:"before"`
	Now        int `json:"now"`
	Conclusion int `json:"conclusion"`
}

type AllClassroomReport struct {
	ClassName string                             `json:"class_name"`
	DOB       string                             `json:"dob"`
	Age       int                                `json:"age"`
	Class     float64                            `json:"class"`
	Topics    map[string]AllClassroomTopicStatus `json:"topics"`
}

type GetReportOverviewAllClassroomResponse struct {
	Reports []AllClassroomReport `json:"reports"`
}
