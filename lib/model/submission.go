package model

type Submission struct {
	ID        string
	ProblemID string `form:"id" binding:"required"`
	Code      string `form:"code" binding:"required"`
}

type SubmissionResult struct {
	stage  string
	status string
}
