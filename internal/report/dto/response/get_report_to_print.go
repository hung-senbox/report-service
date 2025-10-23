package response

type GetReport2Print struct {
	Before     string `json:"before"`
	Now        string `json:"now"`
	Conclusion string `json:"conclusion"`
}
