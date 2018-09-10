package db

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

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

func CreateProblem(p ProblemData) (string, error) {
	problems := db.C("problems")

	id, err := count("problems")
	if err != nil {
		log.Println(err)
		return "", ErrInternalServer
	}

	rec := bson.M{
		"_id":         id,
		"name":        p.Name,
		"description": p.Description,
		"inputs":      p.In,
		"outputs":     p.Out,
	}
	if err := problems.Insert(rec); err != nil {
		log.Println(err)
		return "", ErrInternalServer
	}

	return id, nil
}

func EditProblem(p Problem) error {
	problems := db.C("problems")

	query := bson.M{"_id": p.ID}
	update := bson.M{
		"_id":         p.ID,
		"name":        p.Name,
		"description": p.Description,
		"inputs":      p.In,
		"outputs":     p.Out,
	}

	_, err := problems.Upsert(
		query,
		bson.M{
			"$set": update,
		},
	)
	if err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	return nil
}

func DeleteProblem(pid ProblemID) error {
	problems := db.C("problems")

	query := bson.M{"_id": pid.ID}
	err := problems.Remove(query)
	if err == nil {
		return nil
	} else if err == mgo.ErrNotFound {
		return ErrNoSuchProblem
	} else {
		log.Println(err)
		return ErrInternalServer
	}
}

func Problems() ([]ProblemSnap, error) {
	problems := db.C("problems")

	query := problems.Find(bson.M{})
	count, err := query.Count()
	if err != nil {
		log.Println(err)
		return nil, ErrInternalServer
	}

	res := make([]ProblemSnap, count)

	var (
		problem ProblemSnap
		i       int
	)

	iter := query.Iter()
	for iter.Next(&problem) {
		res[i] = problem
		i++
	}

	err = iter.Close()
	if err != nil {
		log.Println(err)
	}

	return res, nil
}

func prob(id string, format interface{}) error {
	problems := db.C("problems")

	err := problems.Find(bson.M{
		"_id": id,
	}).One(format)

	if err == mgo.ErrNotFound {
		return ErrNoSuchProblem
	} else if err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	return err
}

func FullProblem(id string) (Problem, error) {
	var p Problem
	err := prob(id, &p)
	return p, err
}

func SnapProblem(id string) (ProblemSnap, error) {
	var p ProblemSnap
	err := prob(id, &p)
	return p, err
}
