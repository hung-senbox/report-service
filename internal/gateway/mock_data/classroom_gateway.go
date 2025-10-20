package mockdata

import dto "report-service/internal/gateway/dto/response"

func FakeAllClassroomAssignTemplate() []dto.GetAllClassroomAssignTemplate {
	return []dto.GetAllClassroomAssignTemplate{
		{
			ClassroomID:      "class_1",
			ClassroomName:    "Class A1",
			ClassroomIconUrl: "",
			AssignTemplates: []dto.AssignTemplate{
				{TeacherID: "5849b67e-83e4-4a64-8d79-00611ef26b54", StudentID: "5e470a1e-6b71-4e06-b2cd-ecb7eada4071"},
				{TeacherID: "5849b67e-83e4-4a64-8d79-00611ef26b54", StudentID: "640987ed-0a1f-432a-b006-ed89bbac9d7c"},
			},
		},
		{
			ClassroomID:      "class_2",
			ClassroomName:    "Class B2",
			ClassroomIconUrl: "",
			AssignTemplates: []dto.AssignTemplate{
				{TeacherID: "33e0292c-9e65-4f27-ab54-da4913928f86", StudentID: "5e470a1e-6b71-4e06-b2cd-ecb7eada4071"},
			},
		},
		{
			ClassroomID:      "class_3",
			ClassroomName:    "Class C3",
			ClassroomIconUrl: "",
			AssignTemplates: []dto.AssignTemplate{
				{TeacherID: "5849b67e-83e4-4a64-8d79-00611ef26b54", StudentID: "640987ed-0a1f-432a-b006-ed89bbac9d7c"},
			},
		},
	}
}
