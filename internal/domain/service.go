package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TrainingService interface {
	GetTrainingsByUser(ctx context.Context, userID uuid.UUID) ([]*Training, error)
	GetTrainingWithExercises(ctx context.Context, trainingID int64) (*Training, error)
	CreateTraining(ctx context.Context, cmd CreateTrainingCmd) (*Training, error)
	UpdateTraining(ctx context.Context, cmd UpdateTrainingCmd) (*Training, error)
	DeleteTraining(ctx context.Context, trainingID int64) error
	AddExerciseToTraining(ctx context.Context, cmd AddExerciseToTrainingCmd) (*TrainedExercise, error)
	UpdateTrainedExercise(ctx context.Context, cmd UpdateTrainedExerciseCmd) (*TrainedExercise, error)
	RemoveExerciseFromTraining(ctx context.Context, trainingID, exerciseID int64) error
	GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*TrainingStats, error)
	CompleteTraining(ctx context.Context, trainingID int64, rating *int32) (*Training, error)
}

type CreateTrainingCmd struct {
	UserID            uuid.UUID
	IsDone            bool
	PlannedDate       time.Time
	ActualDate        *time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	TotalDuration     *time.Duration
	TotalRestTime     *time.Duration
	TotalExerciseTime *time.Duration
	Rating            *int32
}

type UpdateTrainingCmd struct {
	ID                int64
	IsDone            *bool
	PlannedDate       time.Time
	ActualDate        *time.Time
	StartedAt         *time.Time
	FinishedAt        *time.Time
	TotalDuration     *time.Duration
	TotalRestTime     *time.Duration
	TotalExerciseTime *time.Duration
	Rating            *int32
}

type AddExerciseToTrainingCmd struct {
	TrainingID int64
	ExerciseID int64
	Weight     *decimal.Decimal
	Approaches *int32
	Reps       *int32
	Time       *time.Duration
	Doing      *time.Duration
	Rest       *time.Duration
	Notes      *string
}

type UpdateTrainedExerciseCmd struct {
	ID         int64
	Weight     *decimal.Decimal
	Approaches *int32
	Reps       *int32
	Time       *time.Duration
	Doing      *time.Duration
	Rest       *time.Duration
	Notes      *string
}

type ExerciseService interface {
	GetAllExercises(ctx context.Context) ([]*Exercise, error)
	GetExerciseByID(ctx context.Context, id int64) (*Exercise, error)
	GetExercisesByTag(ctx context.Context, tagID int64) ([]*Exercise, error)
	SearchExercises(ctx context.Context, query string, tagID *int64) ([]*Exercise, error)
	GetAllTags(ctx context.Context) ([]*Tag, error)
	GetTagByID(ctx context.Context, id int64) (*Tag, error)
	GetExerciseTags(ctx context.Context, exerciseID int64) ([]*Tag, error)
	GetExercisesByMultipleTags(ctx context.Context, tagIDs []int64) ([]*Exercise, error)
	GetPopularTags(ctx context.Context, limit int) ([]*Tag, error)
}
