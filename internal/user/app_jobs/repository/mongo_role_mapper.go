package repository

import (
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoJobMapper struct{}

func NewMongoJobMapper() *MongoJobMapper {
	return &MongoJobMapper{}
}

func (m MongoJobMapper) MapFromEntity(job *entity.Job) (*MongoJob, error) {
	objectID, err := bson.ObjectIDFromHex(job.Id.String())
	if err != nil {
		return nil, err
	}

	return &MongoJob{
		Id:          objectID,
		Name:        job.Name,
		Permissions: job.Permissions,
	}, nil
}

func (m *MongoJobMapper) MapToEntity(mongoJob *MongoJob) (*entity.Job, error) {
	return &entity.Job{
		Id:          entity.JobId(mongoJob.Id.Hex()),
		Name:        mongoJob.Name,
		Permissions: mongoJob.Permissions,
	}, nil
}
