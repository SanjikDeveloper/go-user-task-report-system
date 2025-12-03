package models

type ReportVM struct {
	Period     string `json:"period"`
	Tasks      []Task `json:"tasks"`
	DoneCount  int    `json:"done_count"`
	OtherCount int    `json:"other_count"`
}
