package models

type UpdateQuestion struct {
	Title      string `json:"title"`
	Problem    string `json:"problem"`
	Examples   string `json:"examples"`
	TimeTaken  string `json:"timeTaken"`
	Status     string `json:"status"`
	Difficulty string `json:"difficulty"`
	Notes      string `json:"notes"`
	Tags       []*Tag `json:"tags"`
}
