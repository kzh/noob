package model

type Submission struct {
	ID        string
	ProblemID string `json:"problem"`
	Code      string `json:"code"`
}

type SubmissionResult struct {
	Stage  string `json:"stage"`
	Status string `json:"status"`
}
