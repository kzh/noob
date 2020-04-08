package model

type ProblemID struct {
	ID string `form:"id" binding:"required"`
}

type ProblemData struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
	In          string `form:"inputs" binding:"required"`
	Out         string `form:"outputs" binding:"required"`
}

type Problem struct {
	ID          string `form:"id" json:"id" bson:"_id"`
	Name        string `form:"name" json:"name" bson:"name" binding:"required"`
	Description string `form:"description" json:"description" bson:"description" binding:"required"`
	In          string `form:"inputs" json:"inputs" bson:"inputs" binding:"required"`
	Out         string `form:"outputs" json:"outputs" bson:"outputs" binding:"required"`
}

type ProblemSnap struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type ProblemIO struct {
	In  string `json:"inputs" bson:"inputs"`
	Out string `json:"outputs" bson:"outputs"`
}
