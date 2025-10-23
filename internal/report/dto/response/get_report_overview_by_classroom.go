package response

type GetReportOverviewByClassroomResponse4Web struct {
	ClassInfo              ClassInfo                  `json:"class_info"`
	OverallClassPercentage float32                    `json:"overall_class_percentage"`
	ClassOverview          []ClassOverviewByClassroom `json:"class_overview"`
}

type ClassOverviewByClassroom struct {
	Teacher                 ClassTeacher              `json:"teacher"`
	Student                 ClassStudent              `json:"student"`
	AverageTopicsPercentage float32                   `json:"average_topics_percentage"`
	Topics                  []AllClassroomTopicStatus `json:"topics"`
}

type ClassInfo struct {
	ClassName    string      `json:"class_name"`
	ClassIconUrl string      `json:"class_icon_url"`
	DOB          string      `json:"dob"`
	Age          string      `json:"age"`
	Class        string      `json:"class"`
	Leader       ClassLeader `json:"leader"`
}

type ClassLeader struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ClassTeacher struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ClassStudent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
