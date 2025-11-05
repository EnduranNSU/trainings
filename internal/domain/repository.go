package domain

import (
	"context"

	"github.com/google/uuid"
)

type TrainingRepository interface {
	// Тренировки
	GetTrainingsByUser(ctx context.Context, userID uuid.UUID) ([]*Training, error)
	GetTrainingWithExercises(ctx context.Context, trainingID int64) (*Training, error)
	CreateTraining(ctx context.Context, training *Training) (*Training, error)
	UpdateTraining(ctx context.Context, training *Training) (*Training, error)
	DeleteTrainingAndExercises(ctx context.Context, trainingID int64) error
	
	// Упражнения в тренировках
	AddExerciseToTraining(ctx context.Context, exercise *TrainedExercise) (*TrainedExercise, error)
	UpdateTrainedExercise(ctx context.Context, exercise *TrainedExercise) (*TrainedExercise, error)
	DeleteExerciseFromTraining(ctx context.Context, exerciseID, trainingID int64) error
	
	// Статистика
	GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*TrainingStats, error)
}

type TrainingStats struct {
	TotalTrainings   int64   `json:"total_trainings"`
	CompletedTrainings int64 `json:"completed_trainings"`
	AverageRating    float64 `json:"average_rating"`
	TotalTime        string  `json:"total_time"`
}


type ExerciseRepository interface {
	// Упражнения
	GetExercisesWithTags(ctx context.Context) ([]*Exercise, error)
	GetExerciseByID(ctx context.Context, id int64) (*Exercise, error)
	GetExercisesByTag(ctx context.Context, tagID int64) ([]*Exercise, error)
	SearchExercises(ctx context.Context, filter ExerciseFilter) ([]*Exercise, error)
	
	// Теги
	GetAllTags(ctx context.Context) ([]*Tag, error)
	GetTagByID(ctx context.Context, id int64) (*Tag, error)
	
	// Связи упражнений с тегами
	GetExerciseTags(ctx context.Context, exerciseID int64) ([]*Tag, error)
}