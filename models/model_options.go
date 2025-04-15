package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModelOption struct {
	SortBy  string `json:"sort_by"`
	SortDir int    `json:"sort_dir"`
	Limit   int64  `json:"limit"`
	Skip    int64  `json:"skip"`
}
type ModelOptionMultipleSorting struct {
	Limit int64  `json:"limit"`
	Skip  int64  `json:"skip"`
	Sort  bson.D `json:"sort"`
}

func (modelOption *ModelOption) GetOption(modelOptions []ModelOption) *options.FindOptions {
	findOptions := options.Find()

	if len(modelOptions) == 0 {

		newModelOption := &ModelOption{
			SortBy:  "sequence",
			SortDir: 1,
		}

		sortMap := bson.D{primitive.E{Key: newModelOption.SortBy, Value: newModelOption.SortDir}}
		findOptions.SetSort(sortMap)

		return findOptions
	}

	for _, opt := range modelOptions {
		if opt.Limit > 0 {
			findOptions.SetLimit(opt.Limit)
		}

		if opt.Skip > 0 {
			findOptions.SetSkip(opt.Skip)
		}

		if opt.SortBy == "" {
			opt.SortBy = "sequence"
		}

		if opt.SortDir == 0 {
			opt.SortDir = 1
		}

		sortMap := bson.D{
			primitive.E{Key: opt.SortBy, Value: opt.SortDir},
			primitive.E{Key: "_id", Value: -1},
		}
		findOptions.SetSort(sortMap)
	}

	return findOptions
}

// Format Option before query
func (modelOption *ModelOption) GetAggregateOption(modelOptions []ModelOption) *ModelOptionMultipleSorting {

	modelOptionsFormated := new(ModelOptionMultipleSorting)
	if len(modelOptions) == 0 {
		modelOptionsFormated.Sort = bson.D{primitive.E{Key: "sequence", Value: 1}}
		return modelOptionsFormated
	}

	if modelOptions[0].Limit > 0 {
		modelOptionsFormated.Limit = modelOptions[0].Limit
	}

	if modelOptions[0].Skip > 0 {
		modelOptionsFormated.Skip = modelOptions[0].Skip
	}

	sorts := bson.D{}

	for _, opt := range modelOptions {

		key := "sequence"
		value := 1

		if opt.SortBy != "" {
			key = opt.SortBy
		}
		if opt.SortDir != 0 {
			value = opt.SortDir
		}
		sorts = append(sorts, primitive.E{Key: key, Value: value})
	}

	modelOptionsFormated.Sort = sorts

	return modelOptionsFormated
}
