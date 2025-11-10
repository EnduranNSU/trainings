package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/EnduranNSU/trainings/internal/adapter/out/postgres/gen"
	"github.com/EnduranNSU/trainings/internal/domain"
	"github.com/EnduranNSU/trainings/internal/logging"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type TrainingRepositoryImpl struct {
	q  *gen.Queries
	db *sql.DB
}

func NewTrainingRepository(db *sql.DB) domain.TrainingRepository {
	return &TrainingRepositoryImpl{
		q:  gen.New(db),
		db: db,
	}
}

func (r *TrainingRepositoryImpl) GetTrainingsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Training, error) {
	trainings, err := r.q.GetTrainingsByUser(ctx, userID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": userID.String(),
		})
		logging.Error(err, "GetTrainingsByUser", jsonData, "failed to get trainings by user")
		return nil, err
	}

	result := make([]*domain.Training, len(trainings))
	for i, t := range trainings {
		result[i] = r.toDomainTraining(t)
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"user_id":       userID.String(),
		"trainings_count": len(result),
	})
	logging.Debug("GetTrainingsByUser", jsonData, "successfully retrieved user trainings")

	return result, nil
}

func (r *TrainingRepositoryImpl) GetTrainingWithExercises(ctx context.Context, trainingID int64) (*domain.Training, error) {
	training, err := r.q.GetTrainingWithExercises(ctx, trainingID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": trainingID,
		})
		logging.Error(err, "GetTrainingWithExercises", jsonData, "failed to get training with exercises")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromJoined(training)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id":   trainingID,
		"exercises_count": len(domainTraining.Exercises),
	})
	logging.Debug("GetTrainingWithExercises", jsonData, "successfully retrieved training with exercises")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) CreateTraining(ctx context.Context, training *domain.Training) (*domain.Training, error) {
	params := gen.CreateTrainingParams{
		UserID:    training.UserID,
		Isdone:    training.IsDone,
		Planned:   training.Planned,
		Done:      null.TimeFromPtr(training.Done).NullTime,
		TotalTime: null.IntFromPtr((*int64)(training.TotalTime)).NullInt64,
		Rating:    null.Int32FromPtr(training.Rating).NullInt32,
	}

	created, err := r.q.CreateTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": training.UserID.String(),
			"planned": training.Planned,
			"is_done": training.IsDone,
		})
		logging.Error(err, "CreateTraining", jsonData, "failed to create training")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromGen(created)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": domainTraining.ID,
		"user_id":     domainTraining.UserID.String(),
	})
	logging.Debug("CreateTraining", jsonData, "successfully created training")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) UpdateTraining(ctx context.Context, training *domain.Training) (*domain.Training, error) {
	params := gen.UpdateTrainingParams{
		Isdone:    training.IsDone,
		Planned:   training.Planned,
		Done:      null.TimeFromPtr(training.Done).NullTime,
		TotalTime: null.IntFromPtr((*int64)(training.TotalTime)).NullInt64,
		Rating:    null.Int32FromPtr(training.Rating).NullInt32,
		ID:        training.ID,
	}

	updated, err := r.q.UpdateTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": training.ID,
			"is_done":     training.IsDone,
			"rating":      training.Rating,
		})
		logging.Error(err, "UpdateTraining", jsonData, "failed to update training")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromGen(updated)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": domainTraining.ID,
	})
	logging.Debug("UpdateTraining", jsonData, "successfully updated training")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) DeleteTrainingAndExercises(ctx context.Context, trainingID int64) error {
	err := r.q.DeleteTrainingAndExercises(ctx, trainingID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": trainingID,
		})
		logging.Error(err, "DeleteTrainingAndExercises", jsonData, "failed to delete training and exercises")
		return err
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": trainingID,
	})
	logging.Debug("DeleteTrainingAndExercises", jsonData, "successfully deleted training and exercises")

	return nil
}

func (r *TrainingRepositoryImpl) AddExerciseToTraining(ctx context.Context, exercise *domain.TrainedExercise) (*domain.TrainedExercise, error) {
	params := gen.AddExerciseToTrainingParams{
		TrainingID: exercise.TrainingID,
		ExerciseID: exercise.ExerciseID,
		Weight:     null.FloatFromPtr(exercise.Weight).NullFloat64,
		Approaches: null.IntFromPtr(exercise.Approaches).NullInt64,
		Reps:       null.IntFromPtr(exercise.Reps).NullInt64,
		Time:       null.TimeFromPtr(exercise.Time).NullTime,
		Notes:      null.StringFromPtr(exercise.Notes).NullString,
	}

	created, err := r.q.AddExerciseToTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": exercise.TrainingID,
			"exercise_id": exercise.ExerciseID,
		})
		logging.Error(err, "AddExerciseToTraining", jsonData, "failed to add exercise to training")
		return nil, err
	}

	domainExercise := r.toDomainTrainedExercise(created)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": domainExercise.ID,
		"training_id":         domainExercise.TrainingID,
		"exercise_id":         domainExercise.ExerciseID,
	})
	logging.Debug("AddExerciseToTraining", jsonData, "successfully added exercise to training")

	return domainExercise, nil
}

func (r *TrainingRepositoryImpl) UpdateTrainedExercise(ctx context.Context, exercise *domain.TrainedExercise) (*domain.TrainedExercise, error) {
	params := gen.UpdateTrainedExerciseParams{
		Weight:     null.FloatFromPtr(exercise.Weight).NullFloat64,
		Approaches: null.IntFromPtr(exercise.Approaches).NullInt64,
		Reps:       null.IntFromPtr(exercise.Reps).NullInt64,
		Time:       null.TimeFromPtr(exercise.Time).NullTime,
		Notes:      null.StringFromPtr(exercise.Notes).NullString,
		ID:         exercise.ID,
	}

	updated, err := r.q.UpdateTrainedExercise(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"trained_exercise_id": exercise.ID,
			"training_id":         exercise.TrainingID,
		})
		logging.Error(err, "UpdateTrainedExercise", jsonData, "failed to update trained exercise")
		return nil, err
	}

	domainExercise := r.toDomainTrainedExercise(updated)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": domainExercise.ID,
	})
	logging.Debug("UpdateTrainedExercise", jsonData, "successfully updated trained exercise")

	return domainExercise, nil
}

func (r *TrainingRepositoryImpl) DeleteExerciseFromTraining(ctx context.Context, exerciseID, trainingID int64) error {
	err := r.q.DeleteExerciseFromTraining(ctx, gen.DeleteExerciseFromTrainingParams{
		ID:         exerciseID,
		TrainingID: trainingID,
	})
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"trained_exercise_id": exerciseID,
			"training_id":         trainingID,
		})
		logging.Error(err, "DeleteExerciseFromTraining", jsonData, "failed to delete exercise from training")
		return err
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": exerciseID,
		"training_id":         trainingID,
	})
	logging.Debug("DeleteExerciseFromTraining", jsonData, "successfully deleted exercise from training")

	return nil
}

func (r *TrainingRepositoryImpl) GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*domain.TrainingStats, error) {
	trainings, err := r.GetTrainingsByUser(ctx, userID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": userID.String(),
		})
		logging.Error(err, "GetUserTrainingStats", jsonData, "failed to get user training stats")
		return nil, err
	}

	stats := &domain.TrainingStats{}
	var totalRating int64
	var completedCount int

	for _, t := range trainings {
		stats.TotalTrainings++
		if t.IsDone {
			completedCount++
			if t.Rating != nil {
				totalRating += int64(*t.Rating)
			}
		}
	}

	stats.CompletedTrainings = int64(completedCount)
	if completedCount > 0 {
		stats.AverageRating = float64(totalRating) / float64(completedCount)
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"user_id":            userID.String(),
		"total_trainings":    stats.TotalTrainings,
		"completed_trainings": stats.CompletedTrainings,
		"average_rating":     stats.AverageRating,
	})
	logging.Debug("GetUserTrainingStats", jsonData, "successfully calculated user training stats")

	return stats, nil
}

// Вспомогательные методы для преобразования данных (остаются без изменений)
func (r *TrainingRepositoryImpl) toDomainTraining(t gen.GetTrainingsByUserRow) *domain.Training {
	training := &domain.Training{
		ID:        t.ID,
		UserID:    t.UserID,
		IsDone:    t.Isdone,
		Planned:   t.Planned,
		Done:      nullTimeFromSQL(t.Done),
		TotalTime: (*time.Duration)(nullIntFromSQL(t.TotalTime)),
		Rating:    nullIntFromSQL32(t.Rating),
	}

	// Преобразование упражнений
	if t.Exercises != nil {
		if exercisesSlice, ok := t.Exercises.([]gen.TrainedExercise); ok && len(exercisesSlice) > 0 {
			exercises := make([]domain.TrainedExercise, len(exercisesSlice))
			for i, ex := range exercisesSlice {
				exercises[i] = domain.TrainedExercise{
					ID:         ex.ID,
					TrainingID: training.ID,
					ExerciseID: ex.ExerciseID,
					Weight:     nullFloatFromSQL(ex.Weight),
					Approaches: nullIntFromSQL(ex.Approaches),
					Reps:       nullIntFromSQL(ex.Reps),
					Time:       nullTimeFromSQL(ex.Time),
					Notes:      nullStringFromSQL(ex.Notes),
				}
			}
			training.Exercises = exercises
		}
	}

	return training
}

func (r *TrainingRepositoryImpl) toDomainTrainingFromJoined(t gen.GetTrainingWithExercisesRow) *domain.Training {
	training := &domain.Training{
		ID:        t.ID,
		UserID:    t.UserID,
		IsDone:    t.Isdone,
		Planned:   t.Planned,
		Done:      nullTimeFromSQL(t.Done),
		TotalTime: (*time.Duration)(nullIntFromSQL(t.TotalTime)),
		Rating:    nullIntFromSQL32(t.Rating),
	}

	// Преобразование упражнений
	if t.Exercises != nil {
		if exercisesSlice, ok := t.Exercises.([]gen.TrainedExercise); ok && len(exercisesSlice) > 0 {
			exercises := make([]domain.TrainedExercise, len(exercisesSlice))
			for i, ex := range exercisesSlice {
				exercises[i] = domain.TrainedExercise{
					ID:         ex.ID,
					TrainingID: training.ID,
					ExerciseID: ex.ExerciseID,
					Weight:     nullFloatFromSQL(ex.Weight),
					Approaches: nullIntFromSQL(ex.Approaches),
					Reps:       nullIntFromSQL(ex.Reps),
					Time:       nullTimeFromSQL(ex.Time),
					Notes:      nullStringFromSQL(ex.Notes),
				}
			}
			training.Exercises = exercises
		}
	}

	return training
}

func (r *TrainingRepositoryImpl) toDomainTrainingFromGen(t gen.Training) *domain.Training {
	return &domain.Training{
		ID:        t.ID,
		UserID:    t.UserID,
		IsDone:    t.Isdone,
		Planned:   t.Planned,
		Done:      nullTimeFromSQL(t.Done),
		TotalTime: (*time.Duration)(nullIntFromSQL(t.TotalTime)),
		Rating:    nullIntFromSQL32(t.Rating),
	}
}

func (r *TrainingRepositoryImpl) toDomainTrainedExercise(t gen.TrainedExercise) *domain.TrainedExercise {
	return &domain.TrainedExercise{
		ID:         t.ID,
		TrainingID: t.TrainingID,
		ExerciseID: t.ExerciseID,
		Weight:     nullFloatFromSQL(t.Weight),
		Approaches: nullIntFromSQL(t.Approaches),
		Reps:       nullIntFromSQL(t.Reps),
		Time:       nullTimeFromSQL(t.Time),
		Notes:      nullStringFromSQL(t.Notes),
	}
}