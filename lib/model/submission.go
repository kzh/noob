package model

type Submission struct {
	ID   string `form:"id" binding:"required"`
	Code string `form:"code" binding:"required"`
}
