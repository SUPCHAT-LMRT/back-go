package repository

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "user_jobs"
)

type MongoJobRepositoryDeps struct {
	uberdig.In
	Client    *mongo.Client
	JobMapper *MongoJobMapper
}

type MongoJobRepository struct {
	deps MongoJobRepositoryDeps
}

type MongoJob struct {
	Id            bson.ObjectID `bson:"_id"`
	Name          string        `bson:"name"`
	AssignedUsers []string      `bson:"assigned_users"`
	Permissions   uint64        `bson:"permissions"`
}

func NewMongoJobRepository(deps MongoJobRepositoryDeps) JobRepository {
	return &MongoJobRepository{deps: deps}
}

func (r *MongoJobRepository) FindByName(ctx context.Context, name string) (*entity.Job, error) {
	filter := bson.M{"name": name}

	var mongoJob MongoJob
	err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, filter).
		Decode(&mongoJob)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding job by name: %w", err)
	}

	return r.deps.JobMapper.MapToEntity(&mongoJob)
}

func (r *MongoJobRepository) FindById(
	ctx context.Context,
	jobId entity.JobId,
) (*entity.Job, error) {
	objectID, err := bson.ObjectIDFromHex(jobId.String())
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}

	var mongoJob MongoJob
	err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": objectID}).
		Decode(&mongoJob)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding job: %w", err)
	}

	return r.deps.JobMapper.MapToEntity(&mongoJob)
}

func (r *MongoJobRepository) Create(ctx context.Context, job *entity.Job) error {
	objectID := bson.NewObjectID()
	job.Id = entity.JobId(objectID.Hex())

	mongoJob := &MongoJob{
		Id:            objectID,
		Name:          job.Name,
		Permissions:   job.Permissions,
		AssignedUsers: []string{},
	}

	_, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoJob)
	if err != nil {
		return fmt.Errorf("error creating job: %w", err)
	}

	return nil
}

func (r *MongoJobRepository) Delete(ctx context.Context, jobId entity.JobId) error {
	objectID, err := bson.ObjectIDFromHex(jobId.String())
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error deleting job: %w", err)
	}

	return nil
}

func (r *MongoJobRepository) Update(ctx context.Context, job *entity.Job) error {
	objectID, err := bson.ObjectIDFromHex(string(job.Id))
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":        job.Name,
			"permissions": job.Permissions,
		},
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("error updating job: %w", err)
	}

	return nil
}

func (r *MongoJobRepository) FindAll(ctx context.Context) ([]*entity.Job, error) {
	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error finding jobs: %w", err)
	}
	defer cursor.Close(ctx)

	var mongoJobs []MongoJob
	if err := cursor.All(ctx, &mongoJobs); err != nil {
		return nil, fmt.Errorf("error decoding jobs: %w", err)
	}

	var jobs []*entity.Job
	for _, mongoJob := range mongoJobs {
		job, err := r.deps.JobMapper.MapToEntity(&mongoJob)
		if err != nil {
			return nil, fmt.Errorf("error mapping job: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *MongoJobRepository) AssignToUser(
	ctx context.Context,
	jobId entity.JobId,
	userId user_entity.UserId,
) error {
	objectID, err := bson.ObjectIDFromHex(jobId.String())
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{},
		"$addToSet": bson.M{
			"assigned_users": userId,
		},
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("error assigning job: %w", err)
	}

	return nil
}

func (r *MongoJobRepository) UnassignFromUser(
	ctx context.Context,
	jobId entity.JobId,
	userId user_entity.UserId,
) error {
	objectID, err := bson.ObjectIDFromHex(jobId.String())
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"is_assigned": false,
		},
		"$pull": bson.M{
			"assigned_users": userId,
		},
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("error unassigning job: %w", err)
	}

	return nil
}

func (r *MongoJobRepository) EnsureAdminJobExists(ctx context.Context) (*entity.Job, error) {
	const adminRoleName = "Admin"

	// Vérifiez si le rôle Admin existe déjà
	existingRole, err := r.FindByName(ctx, adminRoleName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification du rôle Admin: %w", err)
	}

	if existingRole != nil {
		// Le rôle Admin existe déjà
		return nil, nil
	}

	// Créez le rôle Admin avec les permissions PermissionAdmin
	adminRole := &entity.Job{
		Name:        adminRoleName,
		Permissions: entity.CREATE_INVITATION | entity.DELETE_INVITATION | entity.ASSIGN_JOB | entity.UNASSIGN_JOB | entity.DELETE_JOB | entity.UPDATE_JOB | entity.UPDATE_JOB_PERMISSIONS | entity.VIEW_ADMINISTRATION_PANEL,
	}

	if err = r.Create(ctx, adminRole); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du rôle Admin: %w", err)
	}

	return adminRole, nil
}

func (r *MongoJobRepository) EnsureManagerJobExists(ctx context.Context) (*entity.Job, error) {
	const managerRoleName = "Manager"

	// Vérifiez si le rôle Admin existe déjà
	existingRole, err := r.FindByName(ctx, managerRoleName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification du rôle Admin: %w", err)
	}

	if existingRole != nil {
		// Le rôle Admin existe déjà
		return nil, nil
	}

	// Créez le rôle Admin avec les permissions PermissionAdmin
	adminRole := &entity.Job{
		Name:        managerRoleName,
		Permissions: entity.CREATE_INVITATION | entity.VIEW_ADMINISTRATION_PANEL,
	}

	if err = r.Create(ctx, adminRole); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du rôle Manageur: %w", err)
	}

	return adminRole, nil
}

//nolint:revive
func (r *MongoJobRepository) FindByUserId(
	ctx context.Context,
	userId user_entity.UserId,
) ([]*entity.Job, error) {
	filter := bson.M{"assigned_users": userId.String()}

	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding jobs for user: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []*entity.Job
	for cursor.Next(ctx) {
		var mongoJob MongoJob
		if err := cursor.Decode(&mongoJob); err != nil {
			return nil, fmt.Errorf("error decoding job: %w", err)
		}

		isAssigned := false
		for _, assignedUser := range mongoJob.AssignedUsers {
			if assignedUser == userId.String() {
				isAssigned = true
				break
			}
		}

		job, err := r.deps.JobMapper.MapToEntity(&mongoJob)
		if err != nil {
			return nil, fmt.Errorf("error mapping job: %w", err)
		}
		job.IsAssigned = isAssigned
		jobs = append(jobs, job)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return jobs, nil
}

//nolint:revive
func (r *MongoJobRepository) UserHasPermission(
	ctx context.Context,
	userId string,
	permission uint64,
) (bool, error) {
	filter := bson.M{"assigned_users": userId}

	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("error finding jobs for user: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mongoJob MongoJob
		if err := cursor.Decode(&mongoJob); err != nil {
			return false, fmt.Errorf("error decoding job: %w", err)
		}

		for _, assignedUser := range mongoJob.AssignedUsers {
			if assignedUser == userId && (mongoJob.Permissions&permission) != 0 {
				return true, nil
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return false, fmt.Errorf("cursor error: %w", err)
	}

	return false, nil
}
