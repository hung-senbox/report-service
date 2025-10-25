package response

import "report-service/internal/report/model"

type ReportTranslateResponse struct {
	Topic        Topic                                  `json:"topic"`
	Translations map[string]model.ReportTranslationData `json:"translations" bson:"translations"`
}

type Topic struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}
