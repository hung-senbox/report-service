package response

type AllClassroomTopicStatus struct {
	TopicID           string  `json:"id"`
	TopicTitle        string  `json:"title"`
	TopicMainImageUrl string  `json:"main_image_url"`
	MainPercentage    float32 `json:"main_percentage"`
	MainStatus        float32 `json:"status"`
	Before            float32 `json:"before"`
	Now               float32 `json:"now"`
	Conclusion        float32 `json:"conclusion"`
}

type AllClassroomReport struct {
	ClassName string                    `json:"class_name"`
	DOB       string                    `json:"dob"`
	Age       int                       `json:"age"`
	Class     float64                   `json:"class"`
	Topics    []AllClassroomTopicStatus `json:"topics"`
}

type GetReportOverviewAllClassroomResponse struct {
	Reports []AllClassroomReport `json:"reports"`
}
