package response

import gw_response "report-service/internal/gateway/dto/response"

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

type ClassOverview struct {
	ClassName               string                    `json:"class_name"`
	DOB                     string                    `json:"dob"`
	Age                     int                       `json:"age"`
	Class                   float32                   `json:"class"`
	AverageTopicsPercentage float32                   `json:"average_topics_percentage"`
	Topics                  []AllClassroomTopicStatus `json:"topics"`
}

type GetReportOverviewAllClassroomResponse4Web struct {
	OverallClassesPercentage float32                     `json:"overall_classes_percentage"`
	Classes                  []ClassOverview             `json:"classes"`
	AllTopics                []gw_response.TopicResponse `json:"all_topics"`
}
